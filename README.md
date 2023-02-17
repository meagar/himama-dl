# himama-dl

An unofficial bulk downloader for [HiMama](https://www.himama.com) "Activities" (photos/videos).

This scrapes your kid's activities through the website, and it's likely to be quite brittle.

`himama-dl` will prompt you for your HiMama credentials, then use it to scrape a list of children/account IDs out of HiMama.com.
It will then allow you to select which child to download activities for, and then output each childs activities in a subdirector.

A typical page of contains activities that look like this:

<img width="953" alt="Screen shot of HiMama activities" src="https://user-images.githubusercontent.com/242474/131499068-e0595e19-df17-48ed-9d88-76eba67913fc.png">

`himama-dl` will download each attached photo or video and name it according to its "Added by", "Date", "Title', and a short hash to ensure uniqueness, for example:

```
2021-08-17 - Preschool Room - Look what I'm doing today - a4993b62.mov
2021-08-18 - Preschool Room - Look what I'm doing today - 1e3db443.jpeg
2021-08-18 - Preschool Room - Look what I'm doing today - cec982de.mov
2021-08-19 - Preschool Room - Look what I'm doing today - 58af6832.jpeg
2021-08-19 - Preschool Room - Look what I'm doing today - 8f0231b4.movg
```

## Installation

```
go install github.com/meagar/himama-dl@latest
```

## Usage

1. Run `himama-dl`; it will prompt for your HiMama credentials
2. Select which child to download data for (or press enter if only one child is found)
3. Wait


