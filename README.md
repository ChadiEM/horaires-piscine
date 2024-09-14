# Horaires pour les piscines Montherlant et Auteuil

... en format JSON.

## Request

`curl localhost:8080/api/piscine`

## Response

```json
[
  {"day":"Mardi","opening-hours":[{"open":"07h00","close":"08h30"},{"open":"11h30","close":"13h30"},{"open":"16h45","close":"22h00"}]},
  {"day":"Mercredi","opening-hours":[{"open":"07h00","close":"08h30"},{"open":"11h30","close":"18h00"}]},
  {"day":"Jeudi","opening-hours":[{"open":"07h00","close":"08h30"},{"open":"11h30","close":"13h30"}]},
  {"day":"Vendredi","opening-hours":[{"open":"07h00","close":"08h30"},{"open":"11h30","close":"13h30"}]},
  {"day":"Samedi","opening-hours":[{"open":"07h00","close":"18h00"}]},
  {"day":"Dimanche","opening-hours":[{"open":"08h00","close":"18h00"}]},
  {"day":"Lundi","opening-hours":[]},
  {"day":"Mardi","opening-hours":[{"open":"10h00","close":"22h00"}]},
  {"day":"Mercredi","opening-hours":[{"open":"07h00","close":"18h00"}]},
  {"day":"Jeudi","opening-hours":[{"open":"07h00","close":"18h00"}]},
  {"day":"Vendredi","opening-hours":[{"open":"07h00","close":"18h00"}]},
  {"day":"Samedi","opening-hours":[{"open":"07h00","close":"18h00"}]},
  {"day":"Dimanche","opening-hours":[{"open":"08h00","close":"18h00"}]}
]
```