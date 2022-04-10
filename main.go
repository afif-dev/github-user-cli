package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/schollz/progressbar/v3"
	"github.com/urfave/cli"
)

func main() {

	app := cli.NewApp()
  app.UseShortOptionHandling = true
	app.Name = "Github User CLI"
	app.Version = "1.0.0"
	app.Compiled = time.Now()
	app.Authors = []cli.Author{
		{
		Name:  "Afif-Dev",
		Email: "afif-dev@grr.la",
		},
	}
	app.Copyright = "(c) 2022 Afif-Dev https://github.com/afif-dev"
	app.Usage = "Get user infos"
	app.UsageText = "Fetch Github user infos"

	app.Commands = []cli.Command{
	  {
      Name:  "user",
      Usage: "fetch user infos by username",
      Flags: []cli.Flag{
        cli.BoolFlag{
          Name: "save, s",
          Usage: "save in json file with avatar",
        },
      },
      Action: func(c *cli.Context) error {
        
        user := c.Args().Get(0)

        if user != "" {
          url := "https://api.github.com/users/" + user
          responseBytes := getData(url)

          var objmap map[string]interface{}
          if err := json.Unmarshal(responseBytes, &objmap); err != nil {
              log.Fatalln(err)
          }

          avatar_url := objmap["avatar_url"].(string)
          // fmt.Println(objmap["avatar_url"])

          j, _ := json.MarshalIndent(objmap, "", "    ");
          fmt.Println(string(j))

          if c.Bool("save") {
            // save in json file
            f, _ := os.Create("users/github_" + user + ".json")
            defer f.Close()
            f.WriteString(string(j))

            // save avatar img
            img, _ := os.Create("users/avatars/github_" + user + ".jpg")
            defer img.Close()

            resp, err := http.Get(avatar_url)
            if err != nil {
              log.Fatalln(err)
            }
            defer resp.Body.Close()
            
            if _, err := io.Copy(img, resp.Body); err != nil {
              log.Fatalln(err)
            }
          }
        } else {
          log.Fatalln("username is required")
        }
        
        return nil
      },
	  },
    {
      Name:  "users",
      Usage: "fetch users infos",
      Flags: []cli.Flag{
        cli.Float64Flag{
          Name: "startId, sid",
          Value: 1,
          Usage: "start with user id",
        },
        cli.IntFlag{
          Name: "totalPage, tp",
          Value: 1,
          Usage: "total page, per page is 100 user data (up to 10 pages)",
        },
        cli.BoolFlag{
          Name: "saveAvatar, sa",
          Usage: "save avatars file",
        },
      },
      Action: func(c *cli.Context) error {

        var objFull []map[string]interface{}
        total_page := c.Int("totalPage")
        last_id :=  c.Float64("startId")
        
        if total_page < 1 || total_page > 10 {
          log.Fatalln("Error: total page is up to 10 pages")
        }

        for i := 0; i < total_page; i++ {
          url := "https://api.github.com/users?since=" + fmt.Sprintf("%.0f", last_id) + "&per_page=100"
          // fmt.Println(url)
          responseBytes := getData(url)

          var objmap []map[string]interface{}
          if err := json.Unmarshal(responseBytes, &objmap); err != nil {
              log.Fatal(err)
          }
          
          last_id = objmap[len(objmap)-1]["id"].(float64)

          for idx := range objmap {
            // append user obj into single obj 
            objFull = append(objFull, objmap[idx])
            
            // save users avatar
            if c.Bool("saveAvatar") {
              // save avatar img
              img, _ := os.Create("users/bulk/avatars/github_" + fmt.Sprintf("%s", objmap[idx]["login"]) + ".jpg")
              defer img.Close()
  
              resp, err := http.Get(fmt.Sprintf("%s", objmap[idx]["avatar_url"]))
              if err != nil {
                log.Fatalln(err)
              }
              defer resp.Body.Close()

              bar := progressbar.DefaultBytes(
                  resp.ContentLength,
                  "downloading avatar - " + fmt.Sprintf("%s", objmap[idx]["login"]),
              )
              
              if _, err := io.Copy(io.MultiWriter(img, bar), resp.Body); err != nil {
                log.Fatalln(err)
              }
            }

          }

          time.Sleep(1 * time.Second)
        }
        
        j, _ := json.MarshalIndent(objFull, "", "    ");
        // fmt.Println(string(j))

        // save user infos in json
        f, err := os.Create("users/bulk/github_"+time.Now().Format("2006-01-02")+".json")
        if err != nil {
          log.Fatalln(err)
        }
        defer f.Close()
        f.WriteString(string(j))

        fmt.Printf("\nSuccess -> Data saved in /users/bulk/ folder")

        return nil
      },
	  },
  }

  err := app.Run(os.Args)
  if err != nil {
    log.Fatal(err)
  }
}
