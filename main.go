package main

func main() {
	if err := dbOpen(); err != nil {
		panic(err)
	}
	defer dbClose()
}
