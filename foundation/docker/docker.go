// Package docker provides support for starting and stopping docker containers
// for running tests.
package docker

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"os/exec"
	"regexp"
	"strings"
	"time"
)

// Container tracks information about the docker container started for tests.
type Container struct {
	Name     string
	HostPort string
}

// StartContainer starts the specified container for running tests.
func StartContainer(image, name, port string, dockerArgs, appArgs []string) (Container, error) {
	// When this code is used in tests, each test could be running in it's own
	// process, so there is no way to serialize the call. The idea is to wait
	// for the container to exist if the code fails to start it.

	if err := validateContainerInputs(name, port); err != nil {
		return Container{}, err
	}

	// Check if the container is already running.
	if c, err := exists(name, port); err == nil {
		return c, nil
	}

	// Try to start the container.
	c, err := dockerRun(image, name, port, dockerArgs, appArgs)
	if err == nil {
		return c, nil
	}

	// The docker run failed. Another test process likely has the container
	// name reserved. Wait for the container to become available.
	for i := range 10 {
		time.Sleep(time.Duration(i+1) * 500 * time.Millisecond)

		c, err := exists(name, port)
		if err == nil {
			return c, nil
		}
	}

	return Container{}, fmt.Errorf("could not start or find container %s", name)
}

// StopContainer stops and removes the specified container.
func StopContainer(id string) error {
	if err := exec.Command("docker", "stop", id).Run(); err != nil {
		return fmt.Errorf("could not stop container: %w", err)
	}

	if err := exec.Command("docker", "rm", id, "-v").Run(); err != nil {
		return fmt.Errorf("could not remove container: %w", err)
	}

	return nil
}

// DumpContainerLogs outputs logs from the running docker container.
func DumpContainerLogs(id string) []byte {
	out, err := exec.Command("docker", "logs", id).CombinedOutput()
	if err != nil {
		return nil
	}

	return out
}

// =============================================================================

func dockerRun(image, name, port string, dockerArgs, appArgs []string) (Container, error) {
	if err := validateContainerInputs(name, port); err != nil {
		return Container{}, err
	}

	safeDockerArgs, err := sanitizeArgs(dockerArgs)
	if err != nil {
		return Container{}, fmt.Errorf("invalid docker args: %w", err)
	}

	safeAppArgs, err := sanitizeArgs(appArgs)
	if err != nil {
		return Container{}, fmt.Errorf("invalid application args: %w", err)
	}

	arg := []string{"run", "-P", "-d", "--name", name}
	arg = append(arg, safeDockerArgs...)
	arg = append(arg, image)
	arg = append(arg, safeAppArgs...)

	var buf bytes.Buffer
	// #nosec G204 - arg is populated with validated/sanitized input.
	cmd := exec.Command("docker", arg...)
	cmd.Stdout = &buf
	if errRun := cmd.Run(); errRun != nil {
		return Container{}, fmt.Errorf("could not start container %s: %w", image, errRun)
	}

	id := buf.String()[:12]
	hostIP, hostPort, err := extractIPPort(id, port)
	if err != nil {
		if serr := StopContainer(id); serr != nil {
			err = errors.Join(err, fmt.Errorf("cleanup failed: %v", serr))
		}
		return Container{}, fmt.Errorf("could not extract ip/port: %w", err)
	}

	c := Container{
		Name:     name,
		HostPort: net.JoinHostPort(hostIP, hostPort),
	}

	return c, nil
}

func exists(name, port string) (Container, error) {
	if err := validateContainerInputs(name, port); err != nil {
		return Container{}, err
	}
	hostIP, hostPort, err := extractIPPort(name, port)
	if err != nil {
		return Container{}, errors.New("container not running")
	}

	c := Container{
		Name:     name,
		HostPort: net.JoinHostPort(hostIP, hostPort),
	}

	return c, nil
}

func extractIPPort(name, port string) (hostIP, hostPort string, err error) {
	if err := validateContainerInputs(name, port); err != nil {
		return "", "", err
	}
	tmpl := fmt.Sprintf("[{{range $k,$v := (index .NetworkSettings.Ports \"%s/tcp\")}}{{json $v}}{{end}}]", port)

	var out bytes.Buffer
	// #nosec G204 - arg is populated with validated/sanitized input.
	cmd := exec.Command("docker", "inspect", "-f", tmpl, name)
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return "", "", fmt.Errorf("could not inspect container %s: %w", name, err)
	}

	// When IPv6 is turned on with Docker.
	// Got  [{"HostIp":"0.0.0.0","HostPort":"49190"}{"HostIp":"::","HostPort":"49190"}]
	// Need [{"HostIp":"0.0.0.0","HostPort":"49190"},{"HostIp":"::","HostPort":"49190"}]
	data := strings.ReplaceAll(out.String(), "}{", "},{")

	var docs []struct {
		HostIP   string `json:"HostIp"`
		HostPort string `json:"HostPort"`
	}
	if err := json.Unmarshal([]byte(data), &docs); err != nil {
		return "", "", fmt.Errorf("could not decode json: %w", err)
	}

	for _, doc := range docs {
		if doc.HostIP != "::" {
			// Podman keeps HostIP empty instead of using 0.0.0.0.
			// - https://github.com/containers/podman/issues/17780
			if doc.HostIP == "" {
				return "localhost", doc.HostPort, nil
			}

			return doc.HostIP, doc.HostPort, nil
		}
	}

	return "", "", fmt.Errorf("could not locate ip/port")
}

var (
	validDockerArg     = regexp.MustCompile(`^[a-zA-Z0-9_\-./:=]+$`)
	validContainerName = regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9_.-]{0,127}$`)
	validPort          = regexp.MustCompile(`^\d{1,5}$`)
)

func validateContainerInputs(name, port string) error {
	if err := validateContainerName(name); err != nil {
		return err
	}

	return validatePort(port)
}

func validateContainerName(name string) error {
	if name == "" || !validContainerName.MatchString(name) {
		return fmt.Errorf("invalid container name %q", name)
	}

	return nil
}

func validatePort(port string) error {
	if port == "" || !validPort.MatchString(port) {
		return fmt.Errorf("invalid port %q", port)
	}

	return nil
}

func sanitizeArgs(args []string) ([]string, error) {
	if len(args) == 0 {
		return nil, nil
	}

	sanitized := make([]string, 0, len(args))
	for i, arg := range args {
		if arg == "" {
			return nil, fmt.Errorf("docker argument %d is empty", i)
		}

		if !validDockerArg.MatchString(arg) {
			return nil, fmt.Errorf("docker argument %q is not allowed", arg)
		}

		sanitized = append(sanitized, arg)
	}

	return sanitized, nil
}
