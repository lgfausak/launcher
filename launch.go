package main

import (
	"fmt"
	"os"
	"regexp"

	"github.com/screwdriver-cd/launcher/screwdriver"
	"github.com/urfave/cli"
)

// VERSION gets set by the build script via the LDFLAGS
var VERSION string

type scmURL struct {
	Host   string
	Org    string
	Repo   string
	Branch string
}

func (s scmURL) String() string {
	return fmt.Sprintf("%s:%s/%s#%s", s.Host, s.Org, s.Repo, s.Branch)
}

// e.g. "git@github.com:screwdriver-cd/launch.git#master"
func parseScmURL(url string) (scmURL, error) {
	r := regexp.MustCompile("(.*):(.*)/(.*)#(.*)")
	matched := r.FindAllStringSubmatch(url, -1)
	if matched == nil || len(matched) != 1 || len(matched[0]) != 5 {
		return scmURL{}, fmt.Errorf("Unable to parse SCM URL %v, match: %q", url, matched)
	}

	return scmURL{
		Host:   matched[0][1],
		Org:    matched[0][2],
		Repo:   matched[0][3],
		Branch: matched[0][4],
	}, nil
}

func launch(api screwdriver.API, buildID string) error {
	b, err := api.BuildFromID(buildID)
	if err != nil {
		return fmt.Errorf("fetching build ID %q: %v", buildID, err)
	}

	j, err := api.JobFromID(b.JobID)
	if err != nil {
		return fmt.Errorf("fetching Job ID %q: %v", b.JobID, err)
	}

	p, err := api.PipelineFromID(j.PipelineID)
	if err != nil {
		return fmt.Errorf("fetching Pipeline ID %q: %v", j.PipelineID, err)
	}

	// org, repo, err := parseScmURL(p.ScmURL)
	// err := git.Checkout(p.ScmURL)
	fmt.Println(p)
	return nil
}

func main() {
	app := cli.NewApp()
	app.Name = "launcher"
	app.Usage = "launch Screwdriver jobs"
	if VERSION == "" {
		VERSION = "0.0.0"
	}
	app.Version = VERSION

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "api-url",
			Usage: "set the API URL for Screwrdriver",
		},
		cli.StringFlag{
			Name:  "tokenfile",
			Usage: "set the JWT used for accessing Screwdriver's API",
		},
	}

	app.Action = func(c *cli.Context) error {
		url := c.String("api-url")
		token := "odsjfadfg"
		buildID := c.Args()[0]

		api, err := screwdriver.New(url, token)
		if err != nil {
			fmt.Printf("creating Scredriver API %v: %v", buildID, err)
		}

		if err = launch(api, buildID); err != nil {
			fmt.Printf("Error running launcher: %v", err)
			os.Exit(1)
		}
		return nil
	}
	app.Run(os.Args)
}
