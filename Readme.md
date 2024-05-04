Socrata - is a command line program that use data from Socrata API and return
10 towns that have the highest sales ratio
and 10 towns that have the highest volume of sales
for each year available in the dataset.

This program is solution for home assignment for Software Engineer position.

### Task
> Using the API of [Socrata](https://dev.socrata.com/foundry/data.ct.gov/5mzw-sjtu) 
we want you to access a [dataset](https://data.ct.gov/resource/5mzw-sjtu.json) 
storing real estate transaction data 
and write a Python or Go program that downloads it, parses a JSON file 
and returns 
10 towns that have the highest sales ratio 
and 10 towns that have the highest volume of sales 
for each year available in the dataset.

> Take home part should not be a “production ready” solution, 
however we encourage you to write a clean and readable code like you would do in real life.

### Run 
**Without token**
```
go run src/main.go
```

**With token**
```
SOCRATA_TOKEN=YOUR_TOKEN go run src/main.go
```

**With token using .env file**

Create .env file in project root directory
Add this line:
```
SOCRATA_TOKEN=YOUR_TOKEN
```
Now you can run
```
go run src/main.go
```

### Questions
Here are several questions that I answered on my own to not distract the hiring team.
But they affected on solution and I would have validated them if it had been my real work task. 
That's why I want to demonstrate them.

---

**Q**: I don't understand what it meas "10 towns that have the highest sales ratio". Every town has several transaction every year and every transaction may have different sales ratio. So do we need to take average of them? Or another metric?

**A**: I had actually the same question for "10 towns that have the highest volume of sales".
I decided to use average as the most common and simplest metric to implement. 
And I made that it is easy to change code to use any other metric. 

---

**Q**: The task has no description program UI. What type of UI should we use?

**A**: As it "should not be a “production ready” solution" I decided to use CLI to minimize amount of code that not solving task. 

---
