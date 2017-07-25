package gorl

import (
	"compress/gzip"
	"encoding/json"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/kennygrant/sanitize"
)

type SaveVersion uint64

const (
	SaveVersionGit SaveVersion = iota
)

const MySaveVersion SaveVersion = SaveVersionGit

type Savedata struct {
	Player  *Critter
	St      *State
	Ow      *Overworld
	Version SaveVersion
}

var (
	configdir string
	datadir   string
)

func SaveGame(player *Critter, state *State, overworld *Overworld) error {
	data := Savedata{
		player,
		state,
		overworld,
		MySaveVersion,
	}
	file, err := os.Create(generateFileName(player, state, overworld))
	if err != nil {
		return err
	}
	defer file.Close()
	zipper := gzip.NewWriter(file)
	defer zipper.Close()
	encoder := json.NewEncoder(zipper)
	err = encoder.Encode(data)
	return err
}

func generateFileName(player *Critter, state *State, overworld *Overworld) string {
	return datadir + string(os.PathSeparator) + sanitize.BaseName(player.Name) + "-" + strconv.Itoa(state.Dungeon) + "-" + time.Now().Format("2006-01-02-15-04-05") + ".json.gz"
}

func LoadGame(state *State) (*Critter, *State, *Overworld, error) {
	var data Savedata
	savefiles, err := ioutil.ReadDir(datadir)
	fnchoice := make([]string, 0, len(savefiles))
	for _, filedata := range savefiles {
		if strings.HasSuffix(filedata.Name(), ".json.gz") && !filedata.IsDir() {
			fnchoice = append(fnchoice, filedata.Name())
		}
	}
	file, err := os.Open(datadir + string(os.PathSeparator) + state.Out.Menu("Load which save?", fnchoice))
	if err != nil {
		return nil, nil, nil, err
	}
	defer file.Close()
	unzipper, err := gzip.NewReader(file)
	if err != nil {
		return nil, nil, nil, err
	}
	defer unzipper.Close()
	decoder := json.NewDecoder(unzipper)
	err = decoder.Decode(&data)
	if err != nil {
		return nil, nil, nil, err
	}
	if MySaveVersion != SaveVersionGit && MySaveVersion != data.Version {
		state.Out.Message("Warning: Save versions differ. You might experience glitches.")
	}
	if data.St.Dungeon <= 0 {
		data.St.CurLevel = data.Ow.M
	} else {
		data.St.CurLevel, data.St.Monsters = data.Ow.M.Tiles[data.Ow.SavedPx][data.Ow.SavedPy].OwData.Dungeon.GetDunLevelFromStorage(data.St.Dungeon)
	}
	return data.Player, data.St, data.Ow, nil
}
