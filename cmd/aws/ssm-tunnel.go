package aws

import (
	"bufio"
	"fmt"
	"math/rand"
	"os/exec"
	"shelllib"
	"strconv"
	"strings"
	"time"

	"github.com/urfave/cli/v2"
)

// NewSSMTunnelCommand return "ssm:tunnel" command
func NewSSMTunnelCommand(globalFlags []cli.Flag) *cli.Command {
	return &cli.Command{
		Name:   "ssm:tunnel",
		Usage:  "Open a SSM tunnel to a remote host",
		Action: SSMTunnel,
		Flags: append(globalFlags, []cli.Flag{
			&cli.StringFlag{
				Name:    "service",
				Aliases: []string{"s"},
				Usage:   "Service Type to connect to (example: bastion, frontend, varnish)",
			},
			&cli.StringFlag{
				Name:    "env",
				Aliases: []string{"e"},
				Usage:   "Environment (example: test, stage, prod)",
			},
			&cli.StringFlag{
				Name:    "instance",
				Aliases: []string{"i"},
				Usage:   "Instace ID (example: i-xxxxxxxxxxxxxxxxx)",
			},
			&cli.StringFlag{
				Name:    "host",
				Aliases: []string{"o"},
				Usage:   "Remote host to open tunnel to (example: myexampledb.a1b2c3d4wxyz.us-west-2.rds.amazonaws.com)",
				Value:   "localhost",
			},
			&cli.StringFlag{
				Name:    "port",
				Aliases: []string{"rp"},
				Usage:   "Remote port to open tunnel to (example: 22)",
				Value:   "22",
			},
			&cli.StringFlag{
				Name:    "local-port",
				Aliases: []string{"lp"},
				Usage:   "Local port to bind to serve tunnel (example: 2222)",
				Value:   "2222",
			},
			&cli.StringFlag{
				Name:    "key",
				Aliases: []string{"k"},
				Usage:   "SSH key to use to connect to instance (example: /path/my-key-pair.pem)",
				Value:   "~/.ssh/id_rsa",
			},
			&cli.StringFlag{
				Name:    "username",
				Aliases: []string{"u"},
				Usage:   "SSH username to use to connect to instance (example: ec2-user)",
				Value:   "ec2-user",
				EnvVars: []string{"USER"},
			},
			&cli.StringFlag{
				Name:   "local-port-ssh",
				Usage:  "Local port used for SSM tunnel (example: 9999)",
				Hidden: true,
			},
		}...),
	}
}

// SSMTunnel open a tunnel to a remote host
func SSMTunnel(c *cli.Context) error {
	var err error
	// Select EC2 instances
	err = SSMSelectInstance(c)
	if err != nil {
		return err
	}

	// Check if SSH tunnel is enough
	localPort := c.String("local-port")
	onlySSH := c.String("port") == "22" && c.String("host") == "localhost"
	if onlySSH == true {
		c.Set("local-port-ssh", localPort)
	}

	// Open SSM tunnel to SSH
	_, err = SSMOpenSSHTunnel(c)
	if err != nil {
		return cli.Exit("Error opening SSH tunnel: "+err.Error(), 1)
	}

	// Notify user that now can connect to tunnel
	fmt.Println(fmt.Sprintf("SSH tunnel to remote instance opened on local port: %s", c.String("local-port-ssh")))

	// Check if an additional SSH tunnel is required
	if onlySSH == true {
		err := shelllib.ExecuteCommandForeground("grep", "/dev/null")
		if err != nil {
			return err
		}
		return nil
	}

	// Open tunnel over SSH
	fmt.Println(fmt.Sprintf("Tunnel to remote %s:%s is available on local port: %s", c.String("host"), c.String("port"), c.String("local-port")))
	err = OpenTunnelOverSSH(c)
	if err != nil {
		return cli.Exit("Error opening tunnel over SSH: "+err.Error(), 1)
	}

	return nil
}

// SSMOpenSSHTunnel open a SSH tunnel using SSM session
func SSMOpenSSHTunnel(c *cli.Context) (*exec.Cmd, error) {
	// Get parameters
	profile := c.String("profile")
	region := c.String("region")
	instanceID := c.String("instance")

	// Elaborate local port
	localPort := c.String("local-port-ssh")
	if len(localPort) == 0 {
		maxPort := 65000
		minPort := 50000
		rand.Seed(time.Now().UnixNano())
		localPort = strconv.Itoa(rand.Intn(maxPort-minPort) + minPort)
		c.Set("local-port-ssh", localPort)
	}

	// Build arguments
	args := []string{
		"ssm", "start-session",
		"--profile", profile,
		"--region", region,
		"--target", instanceID,
		"--document-name", "AWS-StartPortForwardingSession",
		"--parameters", "portNumber=22,localPortNumber=" + localPort,
	}

	// Start SSM session
	cmd, cmdReader, _, err := shelllib.ExecuteCommandBackground("aws", args...)
	if err != nil {
		return cmd, err
	}

	// Wait until a valid response is read
	scanner := bufio.NewScanner(cmdReader)
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), "Port "+localPort+" opened for sessionId") {
			break
		}
	}

	return cmd, nil
}

// OpenTunnelOverSSH open a tunnel using an SSH session
func OpenTunnelOverSSH(c *cli.Context) error {
	// Get parameters
	key := c.String("key")
	username := c.String("username")
	localPortSSH := c.String("local-port-ssh")
	localPort := c.String("local-port")
	host := c.String("host")
	remotePort := c.String("port")

	// Build arguments
	args := []string{
		"-i", key, // SSH key
		"-o", "StrictHostKeyChecking=no", //skip host key verification
		"-p", localPortSSH,
		fmt.Sprintf("%s@localhost", username),
		"-L", fmt.Sprintf("%s:%s:%s", localPort, host, remotePort),
		"-T", // be quite
		"-q", // hide warnings
	}

	// Start SSM session
	err := shelllib.ExecuteCommandForeground("ssh", args...)
	if err != nil {
		return err
	}

	return nil
}
