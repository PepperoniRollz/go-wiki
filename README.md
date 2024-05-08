Playing around with webscraping.

Used Go's net/html package to scrape https://en.wikipedia.org/wiki/Special:AllPages which has a listing of all wikipedia articles.  After collecting all articles, it's stored to disk.  It can be then opened from disk and using a concurrent appraoch, will then create an adjacency list of all wikipedia articles mapped to a list of links on that page.
