package main

import (
	"context"
	"database/sql"
	"fmt"
	"html"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/jeremyu25/blog_aggregator/internal/database"
	"github.com/jeremyu25/blog_aggregator/internal/rss"
)

const defaultLayout = "Mon, 02 Jan 2006 15:04:05 +0000"

func handlerAgg(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <time_between_reqs>\n", cmd.Name)
	}
	time_between_reqs, err := time.ParseDuration(cmd.Args[0])
	if err != nil {
		return err
	}
	fmt.Printf("Collecting Feeds Every %s\n", time_between_reqs)
	ticker := time.NewTicker(time_between_reqs)
	for ; ; <-ticker.C {
		scrapeFeeds(s)
	}
}

func handlerAddFeed(s *state, cmd command, user database.User) error {
	if len(cmd.Args) != 2 {
		return fmt.Errorf("usage: %s <name> <url>", cmd.Name)
	}
	feedName := cmd.Args[0]
	feedUrl := cmd.Args[1]
	feed, err := s.dbQueries.CreateFeed(context.Background(),
		database.CreateFeedParams{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Name:      feedName,
			Url:       feedUrl,
			UserID:    user.ID,
		})
	if err != nil {
		return err
	}
	_, err = s.dbQueries.CreateFeedFollow(context.Background(),
		database.CreateFeedFollowParams{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			UserID:    user.ID,
			FeedID:    feed.ID,
		})
	if err != nil {
		return err
	}
	fmt.Printf("New feed has been created:\nID:%s\nName:%s\nURL:%s\nUser ID:%s\nCreated At:%s\nUpdated At: %s\n", feed.ID, feed.Name, feed.Url, feed.UserID, feed.CreatedAt, feed.UpdatedAt)
	return nil
}
func handlerFeeds(s *state, cmd command) error {
	if len(cmd.Args) != 0 {
		return fmt.Errorf("usage: %s", cmd.Name)
	}
	feeds, err := s.dbQueries.GetAllFeeds(context.Background())
	if err != nil {
		return err
	}
	for i, feed := range feeds {
		fmt.Printf("Feed %d:\n- Name: %s\n- URL: %s\n- User Name: %s\n", i+1, feed.Name, feed.Url, feed.Name_2)
	}
	return nil
}

func handlerFollow(s *state, cmd command, user database.User) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <url>", cmd.Name)
	}
	url := cmd.Args[0]
	feed, err := s.dbQueries.GetFeedByURL(context.Background(), url)
	if err != nil {
		return err
	}
	_, err = s.dbQueries.CreateFeedFollow(context.Background(),
		database.CreateFeedFollowParams{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			FeedID:    feed.ID,
			UserID:    user.ID,
		})
	if err != nil {
		return err
	}
	fmt.Printf("User %s is now following feed %s\n", user.Name, feed.Name)
	return nil
}

func handlerFollowing(s *state, cmd command, user database.User) error {
	if len(cmd.Args) != 0 {
		return fmt.Errorf("usage: %s", cmd.Name)
	}
	feedFollows, err := s.dbQueries.GetFeedFollows(context.Background(), user.ID)
	if err != nil {
		return err
	}
	fmt.Printf("Current User %s is following:\n", s.cfg.CurrentUserName)
	for _, feed := range feedFollows {
		fmt.Printf("- Feed: %s\n- Owned By: %s\n", feed.Name, feed.Name_2)
	}
	return nil
}

func handlerUnfollow(s *state, cmd command, user database.User) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <url>", cmd.Name)
	}
	url := cmd.Args[0]
	feed, err := s.dbQueries.GetFeedByURL(context.Background(), url)
	if err != nil {
		fmt.Printf("Feed url does not exist.\n")
		return err
	}
	_, err = s.dbQueries.DeleteFeedFollow(context.Background(),
		database.DeleteFeedFollowParams{
			UserID: user.ID,
			FeedID: feed.ID,
		})
	if err != nil {
		fmt.Printf("User does not follow this feed.\n")
		return err
	}
	fmt.Println("Successfully unfollowed feed.")
	return nil
}

func scrapeFeeds(s *state) error {
	nextFeed, err := s.dbQueries.GetNextFeedToFetch(context.Background())
	if err != nil {
		return err
	}
	if err = s.dbQueries.MarkFeedFetched(context.Background(),
		database.MarkFeedFetchedParams{
			LastFetchedAt: sql.NullTime{
				Time:  time.Now(),
				Valid: true,
			},
			UpdatedAt: time.Now(),
			ID:        nextFeed.ID,
		}); err != nil {
		return err
	}
	rssFeed, err := rss.FetchFeed(context.Background(), nextFeed.Url)
	if err != nil {
		return err
	}
	for _, item := range rssFeed.Channel.Item {
		postParams, err := setPostParams(item, nextFeed.ID)
		if err != nil {
			return err
		}
		_, err = s.dbQueries.CreatePost(context.Background(), postParams)
		if err != nil {
			return err
		}
	}
	return nil
}

func setPostParams(item rss.RSSItem, feedID uuid.UUID) (database.CreatePostParams, error) {
	postParams := database.CreatePostParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Title:     html.UnescapeString(item.Title),
		Url:       item.Link,
		FeedID:    feedID,
	}
	if item.Description == "" {
		postParams.Description = sql.NullString{
			Valid: false,
		}
	} else {
		postParams.Description = sql.NullString{
			String: html.UnescapeString(item.Description),
			Valid:  true,
		}
	}
	if item.PubDate == "" {
		postParams.PublishedAt = sql.NullTime{
			Valid: false,
		}
	} else {
		time, err := time.Parse(defaultLayout, item.PubDate)
		if err != nil {
			postParams.PublishedAt = sql.NullTime{
				Time:  time,
				Valid: false,
			}
		}
		postParams.PublishedAt = sql.NullTime{
			Time:  time,
			Valid: true,
		}
	}
	return postParams, nil
}

func handlerBrowse(s *state, cmd command, user database.User) error {
	var limit int32 = 2
	if len(cmd.Args) > 1 {
		return fmt.Errorf("usage: %s <optional limit>", cmd.Name)
	}
	if len(cmd.Args) == 1 {
		userLimit, err := strconv.ParseInt(cmd.Args[0], 10, 32)
		if err != nil {
			return err
		}
		limit = int32(userLimit)
	}
	posts, err := s.dbQueries.GetPostsForUser(context.Background(),
		database.GetPostsForUserParams{
			UserID: user.ID,
			Limit:  limit,
		})
	if err != nil {
		return err
	}
	for _, post := range posts {
		fmt.Printf("Title: %s\nDescription: %s\nURL: %s\nPublished At: %s\n\n", post.Title, post.Description.String, post.Url, post.PublishedAt.Time)
	}
	return nil
}
