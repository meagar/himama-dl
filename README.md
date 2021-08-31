# himama-dl

A simple scraper for downloading all "Activities" (photos/videos) of your kids from [HiMama](https://www.himama.com/).

This scrapes your kid's activities through the website, and it's likely to be quite brittle.

`himama-dl` will prompt you for your HiMama credentials, then use it to scrape a list of children/account IDs out of HiMama.com.
It will then allow you to select which child to download activities for, and then output each childs activities in a subdirector.

A typical page of contains activities that look like this:

<img width="953" alt="Screen shot of HiMama activities" src="https://user-images.githubusercontent.com/242474/131499068-e0595e19-df17-48ed-9d88-76eba67913fc.png">

`himama-dl` will download each attached photo or video and name it according to its "Added by", "Date", "Title', and a short hash to ensure uniqueness, for example:

```
2021 08 27 - Preschool 2 Room - Look what I m doing today - fef462ff.jpeg
2021 08 30 - Preschool 2 Room - Look what I m doing today - 3ebdeddb.jpeg
2021 08 30 - Preschool 2 Room - Look what I m doing today - 70cc530f.jpeg
2021 08 30 - Preschool 2 Room - Look what I m doing today - adeff2e9.jpeg
2021 08 30 - Preschool 2 Room - Look what I m doing today - d36a1849.jpeg
```

## Installation

```
go install github.com/meagar/himama-dl
```

## Usage

1. Run `himama-dl`; it will prompt for your HiMama credentials
2. Select which child to download data for (or press enter if only one child is found)
3. Wait


