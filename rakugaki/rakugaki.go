package rakugaki

import (
	"math/rand"
	"weissbot/twitter"
)

// GetRakugaki returns a random tweet that contains either #rkgk, rkgk, or らくがき with more than 100 favorites
func GetRakugaki(tweetSearchFunc func(query string) ([]twitter.Tweet, error)) (twitter.Tweet, error) {
	const favoriteThreshold = 100
	hashtagRkgkTweets, err := tweetSearchFunc("%23rkgk")
	if err != nil {
		return twitter.Tweet{}, err
	}

	rkgkTweets, err := tweetSearchFunc("rkgk")
	if err != nil {
		return twitter.Tweet{}, err
	}

	rkgkJPTweets, err := tweetSearchFunc("%e3%82%89%e3%81%8f%e3%81%8c%e3%81%8d")
	if err != nil {
		return twitter.Tweet{}, err
	}

	var combinedTweets []twitter.Tweet
	combinedTweets = append(combinedTweets, hashtagRkgkTweets...)
	combinedTweets = append(combinedTweets, rkgkTweets...)
	combinedTweets = append(combinedTweets, rkgkJPTweets...)
	topTweets := twitter.FindTopTweets(favoriteThreshold, combinedTweets)

	return topTweets[rand.Intn(len(topTweets))], nil
}
