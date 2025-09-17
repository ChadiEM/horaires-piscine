package piscine

import (
	"io"
	"log"
	"maps"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-shiori/dom"
	"golang.org/x/net/html"
)

var replacer = strings.NewReplacer(
	"Monday", "Lundi",
	"Tuesday", "Mardi",
	"Wednesday", "Mercredi",
	"Thursday", "Jeudi",
	"Friday", "Vendredi",
	"Saturday", "Samedi",
	"Sunday", "Dimanche")

func LoadPiscineMap() map[string]string {
	m := make(map[string]string)

	doc, _ := html.Parse(strings.NewReader(OnPage("https://www.paris.fr/lieux/piscines/tous-les-horaires")))
	rows := dom.QuerySelectorAll(doc, "table.places--timetables-content tbody tr")
	for _, item := range rows {
		row := dom.QuerySelector(item, "td:nth-child(1) a")
		href := dom.GetAttribute(row, "href")
		afterLieux := strings.Split(href, "/")[2]
		key := afterLieux[:strings.LastIndex(afterLieux, "-")]
		m[key] = href
	}

	return m
}

func PiscineHandler(piscineMap map[string]string) func(c *gin.Context) {
	return func(c *gin.Context) {
		value, exists := piscineMap[c.Param("piscine")]
		if !exists {
			c.JSON(404, gin.H{"message": "Piscine not found", "availableValues": slices.Collect(maps.Keys(piscineMap))})
			return
		}

		doc, _ := html.Parse(strings.NewReader(OnPage("https://www.paris.fr" + value)))
		rows := dom.QuerySelectorAll(doc, ".places--schedules-regular-content-title .places--schedules-regular-content-row")[0:7]

		var lastIndex int
		openingHours := make([]Availability, len(rows))
		for index, item := range rows {
			weekday := dom.QuerySelector(item, ".places--schedules-regular-content-weekday")
			realWeekday := dom.InnerText(weekday)

			schedule := dom.QuerySelector(item, ".places--schedules-regular-content-exceptional")
			if schedule == nil {
				schedule = dom.QuerySelector(item, ".places--schedules-regular-content-exceptional-sub:not(.smaller)")
			}

			innerText := dom.InnerText(schedule)
			innerTextSplitByColumn := strings.Split(innerText, ":")
			innerTextMinusTitle := innerTextSplitByColumn[len(innerTextSplitByColumn)-1]

			scheduleTrimmed := strings.TrimSpace(
				innerTextMinusTitle,
			)

			if strings.Contains(scheduleTrimmed, "Fermé") || strings.Contains(scheduleTrimmed, "non renseigné") {
				openingHours[index] = Availability{
					Day:          realWeekday,
					OpeningHours: []Opening{},
				}
			} else {
				scheduleLines := strings.Replace(
					scheduleTrimmed,
					" ", "", -1)

				scheduleLinesWithDashFixed := strings.ReplaceAll(
					scheduleLines, "-", "–")

				realSchedule := strings.Split(
					scheduleLinesWithDashFixed,
					"\n")

				openingHoursSlots := make([]Opening, len(realSchedule))
				for i, schedule := range realSchedule {
					split := strings.Split(schedule, "–")
					openingHoursSlots[i] = Opening{
						Open:  split[0],
						Close: split[1],
					}
				}

				openingHours[index] = Availability{
					Day:          realWeekday,
					OpeningHours: openingHoursSlots,
				}
			}

			lastIndex = index
		}

		lastIndex++

		// fill the remaining hours starting from today and extrapolating into the coming week
		weekday := time.Now().Weekday().String()
		todayIndex := getTodayIndex(openingHours, replacer.Replace(weekday))

		var returnedOpeningHours []Availability

		if len(openingHours) < 7 {
			// something is wrong
			returnedOpeningHours = make([]Availability, 0)
		} else {
			returnedOpeningHours = make([]Availability, max(7, len(openingHours)))
		}

		for i := 0; i < len(returnedOpeningHours); i++ {
			source := openingHours[(todayIndex+i)%len(openingHours)]
			returnedOpeningHours[i] = source
		}

		c.JSON(http.StatusOK, returnedOpeningHours)
	}
}

func getTodayIndex(hours []Availability, today string) int {
	for index, item := range hours {
		if today == item.Day {
			return index
		}
	}
	return -1
}

func OnPage(link string) string {
	res, err := http.Get(link)
	if err != nil {
		log.Fatal(err)
	}
	content, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	err = res.Body.Close()
	if err != nil {
		return ""
	}
	return string(content)
}
