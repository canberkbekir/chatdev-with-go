package utils

import "os"

func WriteToFile(fileName string, fullMessageHistory []string) {
	file, err := os.OpenFile(fileName+".md", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		Logger.Fatal().Err(err).Msg("failed to open file")
		return
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {

		}
	}(file)

	for _, message := range fullMessageHistory {
		if _, err := file.WriteString(message + "\n"); err != nil {
			Logger.Fatal().Err(err).Msg("failed to write to file")
			return
		}
	}
}
