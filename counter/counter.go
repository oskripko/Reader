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

func GetNumberOfAscii(paths []os.FileInfo, dirName string) sync.Map {
	jobs := make(chan string)

	wg := new(sync.WaitGroup)
	jobsWG := new(sync.WaitGroup)

	var counters sync.Map

	for i := 0; i <= workersNumber; i++ {
		wg.Add(1)
		go addOccurrencesNumber(jobs, &counters, wg)
	}
	for _, currentPath := range paths {
		jobsWG.Add(1)
		go func(currentPath os.FileInfo) {
			file, err := os.Open(path.Join(dirName, currentPath.Name()))
			if err != nil {
				log.Fatal(err)
			}

			defer file.Close()
			defer jobsWG.Done()

			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				jobs <- scanner.Text()
			}
		}(currentPath)
	}

	go func() {
		jobsWG.Wait()
		close(jobs)
	}()
	wg.Wait()
	return counters
}

func addOccurrencesNumber(jobs <-chan string, counters *sync.Map, wg *sync.WaitGroup, ) {
	defer wg.Done()

	for line := range jobs {
		for _, value := range line {
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
