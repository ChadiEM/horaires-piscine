package main

import (
	"github.com/gin-gonic/gin"
	"github.com/go-shiori/dom"
	"github.com/samber/lo"
	"golang.org/x/net/html"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

var replacer = strings.NewReplacer(
	"Monday", "Lundi",
	"Tuesday", "Mardi",
	"Wednesday", "Mercredi",
	"Thursday", "Jeudi",
	"Friday", "Vendredi",
	"Saturday", "Samedi",
	"Sunday", "Dimanche")

func main() {
	r := gin.Default()
	r.ForwardedByClientIP = true
	r.SetTrustedProxies([]string{"127.0.0.1"})
	r.GET("/api/horaires-piscine", piscineHandler())
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

func piscineHandler() func(c *gin.Context) {
	return func(c *gin.Context) {
		doc, _ := html.Parse(strings.NewReader(OnPage("https://www.paris.fr/lieux/piscine-henry-de-montherlant-2939")))
		rows := dom.QuerySelectorAll(doc, ".places--schedules-regular-content-title .places--schedules-regular-content-row")
		var openingHours [14]Availability
		for index, item := range rows {
			weekday := dom.QuerySelector(item, ".places--schedules-regular-content-weekday")
			realWeekday := dom.InnerText(weekday)

			schedule := dom.QuerySelector(item, ".places--schedules-regular-content-exceptional")
			if schedule == nil {
				schedule = dom.QuerySelector(item, ".places--schedules-regular-content-exceptional-sub:not(.smaller)")
			}

			scheduleTrimmed := strings.TrimSpace(
				dom.InnerText(schedule),
			)

			if strings.Contains(scheduleTrimmed, "Fermé") {
				openingHours[index] = Availability{
					Day:          realWeekday,
					OpeningHours: []Opening{},
				}
			} else {
				scheduleLines := strings.Replace(
					scheduleTrimmed,
					" ", "", -1)
				realSchedule := strings.Split(
					scheduleLines,
					"\n")

				openingHoursSlots := lo.Map[string, Opening](realSchedule, func(x string, _ int) Opening {
					split := strings.Split(x, "–")
					return Opening{
						Open:  split[0],
						Close: split[1],
					}
				})

				openingHours[index] = Availability{
					Day:          realWeekday,
					OpeningHours: openingHoursSlots,
				}
			}
		}

		weekday := time.Now().Weekday().String()
		todayIndex := getTodayIndex(openingHours, replacer.Replace(weekday))

		c.JSON(http.StatusOK, openingHours[todayIndex:])
	}
}

func getTodayIndex(hours [14]Availability, today string) int {
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
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
	return string(content)
}

type Availability struct {
	Day          string    `json:"day"`
	OpeningHours []Opening `json:"opening-hours"`
}

type Opening struct {
	Open  string `json:"open"`
	Close string `json:"close"`
}
