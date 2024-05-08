Playing around with webscraping.

Used Go's net/html package to scrape https://en.wikipedia.org/wiki/Special:AllPages which has a listing of all wikipedia articles.  After collecting all articles, it's stored to disk.  It can be then opened from disk and using a concurrent appraoch, will then create an adjacency list of all wikipedia articles mapped to a list of links on that page.

The amount of articles crawled through on first pass is ~17 million, which is a bit higher than the Googled amount of ~16.5 million.  Possibly due to redirects or duplicate articles found.
