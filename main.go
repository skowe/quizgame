package main

import (
    "strings"
    "fmt"
    "os"
    "flag"
    "encoding/csv"
    "time"
)

func main(){
    // Define the -csv and -timer flags
    csvFilename := flag.String("csv", "quizzes.csv",
    "A csv file that follows the question, answer format")
    timer := flag.Int("timer", 30, "The duration of the quiz in seconds")

    flag.Parse()

    q, err := NewQuiz(*csvFilename, *timer)
    if err != nil{
        os.Exit(1)
    }
    quizloop:
        for i, problem := range q.problems {
            fmt.Println(problem.question,"=")
            go q.Answer(i)
            select{
            case <-q.answerCh:
                break
            case <-q.time:
                fmt.Println("Time is up!")
                break quizloop
                //return
            }
        }
    q.Rep()
}

type Quiz struct{
    // A structure that cointains most of the needed helper functionality
    // Holds the question/answer pairs
    // Provides the channel for timer to keep track of countdown on a separate thread
    // Provides the channel for answers to register user input on a separate thread
    // Keeps track of all the correct answers
    problems []problem
    time <-chan time.Time
    answerCh chan string
    count int
}

func NewQuiz(file string, timer int) (*Quiz, error){
    // Parses the csv file to create a slice of problems
    // Create a channel for answers
    // Creates a new timer at the end to begin the quizz

    fl, err:= os.Open(file)
    defer fl.Close()
    if err != nil{
        fmt.Println("Failed to open the csv file:",file)
        return nil, err
    }
    csvReader:= csv.NewReader(fl)
    pairs, err := csvReader.ReadAll()
    if err != nil{
        fmt.Println("Failed to parse the contents of the csv file:", file)
        return nil, err
    }
    prbs := make([]problem, len(pairs))
    for i, pair:= range pairs{
        prbs[i].question = pair[0]
        prbs[i].answer = strings.TrimSpace(pair[1])
    }

    q:= Quiz{
        count: 0,
        problems: prbs,
        answerCh: make(chan string),
        time: time.NewTimer(time.Duration(timer)*time.Second).C,
    }

    return &q, nil
}

func (q Quiz) Rep(){
    // Provides a report on the end of the quiz
    fmt.Printf("Answered %d pf %d\n", q.count, len(q.problems))
}

func (q *Quiz) Answer(index int){
    // Recieve the answer from user input
    // Increment the coun if the answers is correct
    // Send the answer down the answer channel
    var ans string
    fmt.Scanf("%s", &ans)
    if ans == q.problems[index].answer{
        q.count++
    }
    q.answerCh <- ans
}

type problem struct{
    question, answer string
}
