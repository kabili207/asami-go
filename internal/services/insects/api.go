package insects

import (
	"fmt"
	"math/rand"
	"net/http"
	"reflect"
	"regexp"
	"strings"
	"time"

	goquery "github.com/PuerkitoBio/goquery"
)

const apiBaseUrl string = "https://butterfly-conservation.org/"
const mothBaseUrl string = "moths/"
const butterflyBaseUrl string = "uk-butterflies/"

type InsectType string

const (
	Moth      InsectType = "moth"
	Butterfly InsectType = "butterfly"
)

var mappedFields = map[string]string{
	"size and family":                    "SizeAndFamily",
	"conservation status":                "ConservationStatus",
	"caterpillar food plants":            "CaterpillarFood",
	"particular caterpillar food plants": "CaterpillarFood",
	"habitat":                            "Habitat",
	"flight season":                      "FlightSeason",
	"distribution":                       "Distribution",
}

func GetRandomInsect(insectType InsectType) (*Insect, error) {
	insects, err := GetInsectList(insectType)
	if err != nil {
		return nil, err
	}
	rand.Seed(time.Now().Unix())
	return GetInsect(insectType, insects[rand.Intn(len(insects))].Key)
}

func GetInsectList(insectType InsectType) ([]InsectSummary, error) {
	var baseUrl string

	if insectType == Moth {
		baseUrl = mothBaseUrl
	} else {
		baseUrl = butterflyBaseUrl
	}

	doc, err := getPage(baseUrl + "a-to-z")
	if err != nil {
		return nil, err
	}
	mothNodes := doc.Find(".atoz li span.field-content a")
	insects := make([]InsectSummary, len(mothNodes.Nodes))
	mothNodes.Each(func(i int, s *goquery.Selection) {
		link, _ := s.Attr("href")
		insects[i] = InsectSummary{
			Type: Moth,
			Name: strings.TrimSpace(s.Text()),
			Key:  strings.TrimPrefix(link, "/"+mothBaseUrl),
		}
	})
	return insects, err
}

func GetInsect(insectType InsectType, key string) (*Insect, error) {

	var baseUrl string

	if insectType == Moth {
		baseUrl = mothBaseUrl
	} else {
		baseUrl = butterflyBaseUrl
	}
	doc, err := getPage(baseUrl + key)

	if err != nil {
		return nil, err
	}

	insect := Insect{
		ID:             key,
		Name:           strings.TrimSpace(doc.Find("div#block-pagetitle h1").Text()),
		ScientificName: strings.TrimSpace(doc.Find("div#block-butterflytagsblock p.sub-heading").Text()),
	}
	parseDataList(doc.Find("#block-butterfly-content div.col:first-of-type .colpad > div"), &insect)

	images := doc.Find("div.feature-item")
	pictures := make([]PictureInfo, len(images.Nodes))
	images.Each(func(i int, q *goquery.Selection) {
		url, _ := q.Find("img").Attr("src")
		pictures[i] = PictureInfo{
			Url:         apiBaseUrl + strings.TrimPrefix(url, "/"),
			Description: q.Find(".feature-text h3").Text(),
			Credit:      q.Find(".feature-text p.gallery-credit").Text(),
		}
	})

	insect.Pictures = pictures
	return &insect, nil
}

func parseDataList(sel *goquery.Selection, insect *Insect) {

	currSection := "Description"
	sel = sel.Children()
	stringVals := []string{}

	sel.Each(func(i int, q *goquery.Selection) {
		switch tagType := q.Nodes[0].Data; tagType {
		case "h4":
			setField(insect, currSection, stringVals)
			stringVals = []string{}
			currSection = mappedFields[strings.ToLower(q.Text())]
		case "p":
			stringVals = append(stringVals, strings.TrimSpace(q.Text()))
		case "ul":
			stringVals = q.Find("li").Map(func(i int, q *goquery.Selection) string {
				return strings.TrimSpace(q.Text())
			})
		}
	})

	setField(insect, currSection, stringVals)
}

func setField(insect *Insect, currentSection string, values []string) {
	if currentSection == "Description" {
		insect.Description = strings.Join(values, "\n\n")
	} else if currentSection == "SizeAndFamily" {
		parseSizeAndFamily(insect, values)
	} else if currentSection == "ConservationStatus" {
		parseConservationStatus(insect, values)
	} else if currentSection == "Distribution" {
		parseDistributionStatus(insect, values)
		setFieldReflect(insect, "DistributionRaw", values)
	} else {
		setFieldReflect(insect, currentSection, values)
	}
}

func normalizeCharacters(value string) string {
	spaceRE := regexp.MustCompile(`\p{Zs}`)
	dashRE := regexp.MustCompile(`\p{Pd}`)

	value = spaceRE.ReplaceAllString(value, " ")
	return dashRE.ReplaceAllString(value, "-")
}

func parseSizeAndFamily(insect *Insect, values []string) {
	familyRE := regexp.MustCompile(`(?i)Family\s+-\s+([\w\(\)\s,-]+)`)
	sizeRE := regexp.MustCompile(`(?i)([\w\s/]+)[\s-]Sized`)
	wingspanRE := regexp.MustCompile(`(?i)wingspan(?:\srange)\s+-\s+([\w\s-]+)`)

	for _, v := range values {
		v = normalizeCharacters(v)
		if familyRE.MatchString(v) {
			insect.Family = familyRE.FindStringSubmatch(v)[1]
		} else if sizeRE.MatchString(v) {
			insect.Size = sizeRE.FindStringSubmatch(v)[1]
		} else if wingspanRE.MatchString(v) {
			insect.Wingspan = wingspanRE.FindStringSubmatch(v)[1]
		}
	}
}

func parseConservationStatus(insect *Insect, values []string) {
	ukBapRE := regexp.MustCompile(`(?i)UK BAP:\s+([\w\(\)\s,-]+)`)
	sizeRE := regexp.MustCompile(`(?i)([\w\s/]+)[\s-]Sized`)

	for _, v := range values {
		v = normalizeCharacters(v)
		if ukBapRE.MatchString(v) {
			insect.ConservationStatus.UK_BAP = ukBapRE.FindStringSubmatch(v)[1]
		} else if sizeRE.MatchString(v) {
			insect.Size = sizeRE.FindStringSubmatch(v)[1]
		} else {
			insect.ConservationStatus.General = strings.TrimSpace(insect.ConservationStatus.General + "\r\n" + v)
		}
	}
}

func parseDistributionStatus(insect *Insect, values []string) {
	countriesRE := regexp.MustCompile(`(?i)Countries\s*[:-]\s+([\w\(\)\s,-]+)`)
	localitiesRE := regexp.MustCompile(`(?i)Localities\s*[:-]\s+([\w\(\)\s,-]+)`)
	trendRE := regexp.MustCompile(`(?i)Distribution Trend Since 1970['’]?s\s+=\s+([\w\s:]+)`)

	for _, v := range values {
		v = normalizeCharacters(v)
		if countriesRE.MatchString(v) {
			countries := countriesRE.FindStringSubmatch(v)[1]
			insect.Distribution.Countries = processList(countries)
		} else if localitiesRE.MatchString(v) {
			localities := localitiesRE.FindStringSubmatch(v)[1]
			insect.Distribution.Localities = processList(localities)
		} else if trendRE.MatchString(v) {
			insect.Distribution.TrendSince1970s = trendRE.FindStringSubmatch(v)[1]
		} else {
			insect.Distribution.General = strings.TrimSpace(insect.Distribution.General + "\r\n" + v)
		}
	}
}

func setFieldReflect(item interface{}, field string, value interface{}) error {
	v := reflect.ValueOf(item).Elem()
	if !v.CanAddr() {
		return fmt.Errorf("cannot assign to the item passed, item must be a pointer in order to assign")
	}
	fv := v.FieldByName(field)

	if fv.IsValid() && fv.CanSet() {
		newVal := reflect.ValueOf(value)
		if newVal.Kind() == reflect.Slice && fv.Kind() == reflect.String {
			newVal = reflect.ValueOf(strings.Join(value.([]string), "\r\n"))
		}
		fv.Set(newVal)

	}
	return nil
}

func processList(rawString string) []string {
	if strings.Contains(rawString, " and ") {
		rawString = strings.ReplaceAll(rawString, " and ", ",")
	}
	data := strings.Split(rawString, ",")
	for i, v := range data {
		data[i] = strings.TrimSpace(v)
	}
	return data
}

func getPage(url string) (*goquery.Document, error) {
	// Request the HTML page.
	res, err := http.Get(apiBaseUrl + url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return nil, err
	}

	// Load the HTML document
	return goquery.NewDocumentFromReader(res.Body)
}
