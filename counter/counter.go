package counter

import (
	"bufio"
	"log"
	"os"
	"path"
	"sync"
	"unicode"
)

const (
	workersNumber = 5
)

type job struct {
	text string
}

func GetNumberOfAscii(paths []os.FileInfo, dirName string) sync.Map {
	jobs := make(chan job, 5)

	wg := new(sync.WaitGroup)
	jobsWG := new(sync.WaitGroup)

	var counters sync.Map

	for i := 0; i <= workersNumber; i++ {
		wg.Add(1)
		go addOccurrencesNumber(jobs, &counters, wg)
	}
	for i := range paths {
		jobsWG.Add(1)
		go func(i int) {
			file, err := os.Open(path.Join(dirName, paths[i].Name()))
			if err != nil {
				log.Fatal(err)
			}

			defer file.Close()
			defer jobsWG.Done()

			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				jobs <- job{text: scanner.Text()}
			}
		}(i)
	}

	go func() {
		jobsWG.Wait()
		close(jobs)
	}()
	wg.Wait()
	return counters
}

func addOccurrencesNumber(jobs <-chan job, counters *sync.Map, wg *sync.WaitGroup, ) {
	defer wg.Done()

	for line := range jobs {
		for _, value := range line.text {
			if unicode.IsSpace(value) {
				continue
			}
			//Возможно стоит предзагрузить все элементы для того, чтобы убрать лишнюю операцию чтения
			currentCounter, exists := counters.LoadOrStore(string(value), 1)
			if exists {
				newCounter := currentCounter.(int)
				newCounter++
				counters.Store(string(value), newCounter)
			}
		}
	}
}
