package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

func main() {
	username := flag.String("username", "tbxark", "GitHub username")
	token := flag.String("token", "", "GitHub token")
	output := flag.String("output", fmt.Sprintf("starred-%s-%d.csv", *username, time.Now().Unix()), "Output file")
	flag.Parse()
	if *username == "" || *token == "" {
		flag.PrintDefaults()
		return
	}
	file, err := os.Create(*output)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}
	_, _ = file.WriteString(strings.Join([]string{
		"name",
		"full_name",
		"html_url",
		"description",
		"stargazers_count",
		"forks_count",
		"topics",
		"language",
		"created_at",
		"updated_at",
	}, ","))
	_, _ = file.WriteString("\n")
	encodeCSVComment := func(s string) string {
		if !strings.ContainsRune(s, ',') {
			return s
		}
		return fmt.Sprintf(`"%s"`, strings.ReplaceAll(s, `"`, `""`))
	}
	_, err = fetchAllStarred(*username, *token, func(repos []Repo) error {
		for _, repo := range repos {
			_, _ = file.WriteString(strings.Join([]string{
				repo.Name,
				repo.FullName,
				repo.HtmlUrl,
				encodeCSVComment(repo.Description),
				fmt.Sprintf("%d", repo.StargazersCount),
				fmt.Sprintf("%d", repo.ForksCount),
				encodeCSVComment(strings.Join(repo.Topics, ",")),
				repo.Language,
				repo.CreatedAt.Format(time.RFC3339),
				repo.UpdatedAt.Format(time.RFC3339),
			}, ","))
			_, _ = file.WriteString("\n")
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
	_ = file.Sync()
}
