package main


import (
    "path/filepath"
    "os"
    //"fmt"
    //"log"
    "strings"
    "time"
    "sort"
    
    //go get "github.com/robfig/cron"
    "github.com/robfig/cron"
)

//BASE

func cronSet() {
    
    if len( G_CRONSHORTTIME ) > 0 || len( G_CRONLONGTIME ) > 0 {
        c := cron.New()
        if len( G_CRONSHORTTIME ) > 0 {
            c.AddFunc( G_CRONSHORTTIME, cronShortRun)
            showInfo( "CRON-SHORT-ACTIVED: " + G_CRONSHORTTIME )
        }
        if len( G_CRONLONGTIME ) > 0 {
            c.AddFunc( G_CRONLONGTIME, cronLongRun)
            showInfo( "CRON-LONG-ACTIVED: " + G_CRONLONGTIME )
        }
        c.Start()
    }
    
}

//SHORT CRON

func cronShortRun(){
    showInfo( "CRON-SHORT-STARTED: " + dateGetNow() )
    fileRemove( G_CRONSHORTTIME_FILE )
    fileAppendLine( G_CRONSHORTTIME_FILE, "Cron SHORT START: " + dateGetNow() + "\n" )
    
    //add new files
    cronAddNewFiles()
    
    //Clean Not Existant Files
    cronCleanFiles()
    
    //Clean Temp Folders
    cronCleanTempFolders()
    
    //Identify new files
    cronIdentNewFiles()
    
    fileAppendLine( G_CRONSHORTTIME_FILE, "Cron SHORT ENDED: " + dateGetNow() )
    showInfo( "CRON-SHORT-ENDED: " + dateGetNow() )
}

//LONG CRON

func cronLongRun(){
    showInfo( "CRON-LONG-STARTED: " + dateGetNow() )
    fileRemove( G_CRONLONGTIME_FILE )
    fileAppendLine( G_CRONLONGTIME_FILE, "Cron LONG START: " + dateGetNow() + "\n" )
    
    //Clean playing
    cronCleanPlaying()
    
    //Link same images
    cronImagesLink()
    
    //Clean IP whitelist old IPs
    cronCleanWhitelist()
    cronCleanBans()
    
    //Clean sessions logins old now-1year
    cronCleanSessions()
    
    //Clean downloaded media duply
    cronCleanMediaDuplys()
    
    //Clean downloaded media not identified
    cronCleanMediaNotIdent()
    
    //Clean Mediainfo duplys
    cronCleanMediaInfoDuplys()
    
    //Compressed Files process
    cronUncompressFiles()
    
    //Free space on low
    cronCleanDisk()
    
    //Clean low size old directories/files
    cronCleanDirsLowSize()
    
    //Complete mediainfo images
    cronMediaInfoCompleteImgs()
    
    //Get Own Searchs and pass to webscrapp downloader
    
    //Add mediainfo images
    
    //Sites Scrap
    
    //LiveTV Clean && Update (1week)
    cronLiveTV()
    
    fileAppendLine( G_CRONLONGTIME_FILE, "Cron LONG ENDED: " + dateGetNow() )
    showInfo( "CRON-LONG-ENDED: " + dateGetNow() )
}

//ACTIONS

//Scan New Files

func cronAddNewFiles(){
    //G_DOWNLOADS_FOLDER
    fileAppendLine( G_CRONSHORTTIME_FILE, "::Adding New Files: " + dateGetNow() )
    nowpath := pathAbs(G_DOWNLOADS_FOLDER)
    fileAppendLine( G_CRONSHORTTIME_FILE, "::Folder: " + nowpath )
    fileAppendLine( G_CRONSHORTTIME_FILE, "" )
    err := filepath.Walk(nowpath,
        func(path string, info os.FileInfo, err error) error {
            if err != nil {
                return err
            }
            if err == nil && fileExist( path ) && checkIsFile( path ) && sliceInString( path, G_DOWNLOADS_FOLDER_EXC ) == false {
                showInfo( "CRON-SHORT-ADDNEWFILES: " + path )
                if sqlite_checkMediaFile( path ) == false && checkMimeVideo( path ) {
                    showInfo( "CRON-SHORT-ADDNEWFILES-ADDED: " + path )
                    sqlite_media_insert( path, "0" )
                    fileAppendText( G_CRONSHORTTIME_FILE, "+" )
                }
            }
            return nil
    })
    fileAppendLine( G_CRONSHORTTIME_FILE, ":: END Adding New Files: " + dateGetNow() )
    if err != nil {
        showInfoError(err)
    }
}

//Clean Not Found Files

func cronCleanFiles(){
    //G_DOWNLOADS_FOLDER
    fileAppendLine( G_CRONLONGTIME_FILE, "::Clean Removed Files: " + dateGetNow() )
    fileAppendLine( G_CRONLONGTIME_FILE, "" )
    filemedia := sqlite_getMediaAll()
    fileAppendLine( G_CRONLONGTIME_FILE, "::Filenum: " + intToStr( len( filemedia ) ) )
    for _, media := range filemedia {
        if media[ "file" ] != "" && fileExist( media[ "file" ] ) == false {
            sqlite_media_delete( media[ "idmedia" ] )
            fileAppendText( G_CRONLONGTIME_FILE, "-" )
        } else {
            fileAppendText( G_CRONLONGTIME_FILE, "=" )
        }
    }
    fileAppendLine( G_CRONLONGTIME_FILE, ":: END Clean Removed Files: " + dateGetNow() )
}

//Identify new files

func cronIdentNewFiles(){
    //G_DOWNLOADS_FOLDER
    nowpath := pathAbs(G_DOWNLOADS_FOLDER)
    fileAppendLine( G_CRONSHORTTIME_FILE, "::Identify New Files: " + dateGetNow() )
    fileAppendLine( G_CRONSHORTTIME_FILE, "" )
    filemedia := sqlite_getMediaIdentNow( 50 )
    fileAppendLine( G_CRONSHORTTIME_FILE, "::Filenum: " + intToStr( len( filemedia ) ) )
    for _, media := range filemedia {
        if media[ "file" ] != "" && fileExist( media[ "file" ] ) {
            idmediainfo := identMedia( media[ "idmedia" ], "mydb", "", "", "", "" )
            if idmediainfo == "" {
                idmediainfo = identMedia( media[ "idmedia" ], G_CRONSCRAPPER, "", "", "", "" )
            }
            filesub := strings.Replace(media[ "file" ], nowpath, "", -1)
            if idmediainfo != "" {
                midata := sqlite_getMediaInfoID(idmediainfo)
                if len( midata ) > 0 {
                    fileAppendLine( G_CRONSHORTTIME_FILE, "::File Title Added: " + midata[0]["title"] + " (" + midata[0]["year"] + ") -> " + filesub )
                } else {
                    fileAppendLine( G_CRONSHORTTIME_FILE, "::File Title Added ADDED BUT NOT DATA TO idmediainfo: " + idmediainfo )
                }
            } else {
                fileAppendLine( G_CRONSHORTTIME_FILE, "::File Failed: " + filesub )
            }
        }
    }
    fileAppendLine( G_CRONLONGTIME_FILE, ":: END Clean Removed Files: " + dateGetNow() )
}

//Clean TMP folder

func cronCleanTempFolders(){
    bf := pathAbs(G_TMP_FOLDER)
    folders := getFolders( bf )
    fileAppendLine( G_CRONSHORTTIME_FILE, "::Removing TMP Folders: " + dateGetNow() )
    for _, f := range folders {
        if delTree( f ) {
            fileAppendLine( G_CRONSHORTTIME_FILE, "::Folder Removed: " + f )
        } else {
            fileAppendLine( G_CRONSHORTTIME_FILE, "::Folder Remove FAILED: " + f )
        }
    }
    fileAppendLine( G_CRONSHORTTIME_FILE, "::END Removing TMP Folders: " + dateGetNow() )
}

//Clean old playing media

func cronCleanPlaying(){
    //G_DOWNLOADS_FOLDER
    fileAppendLine( G_CRONLONGTIME_FILE, "::Clean Playing Files: " + dateGetNow() )
    fileAppendLine( G_CRONLONGTIME_FILE, "" )
    sqlite_playing_clean()
    fileAppendLine( G_CRONLONGTIME_FILE, ":: END Clean Playing Files: " + dateGetNow() )
}

//Clean old whitelist ip

func cronCleanWhitelist(){
    //G_DOWNLOADS_FOLDER
    fileAppendLine( G_CRONLONGTIME_FILE, "::Clean Whitelist IP: " + dateGetNow() )
    fileAppendLine( G_CRONLONGTIME_FILE, "" )
    sqlite_whitelist_clean()
    fileAppendLine( G_CRONLONGTIME_FILE, ":: END Clean  Whitelist IP: " + dateGetNow() )
}

//Clean old bans ip

func cronCleanBans(){
    //G_DOWNLOADS_FOLDER
    fileAppendLine( G_CRONLONGTIME_FILE, "::Clean BANS IP: " + dateGetNow() )
    fileAppendLine( G_CRONLONGTIME_FILE, "" )
    sqlite_bans_clean()
    fileAppendLine( G_CRONLONGTIME_FILE, ":: END Clean BANS IP: " + dateGetNow() )
}

//Clean old Sessions ip

func cronCleanSessions(){
    //G_DOWNLOADS_FOLDER
    fileAppendLine( G_CRONLONGTIME_FILE, "::Clean Sessions IP: " + dateGetNow() )
    fileAppendLine( G_CRONLONGTIME_FILE, "" )
    sqlite_sessions_clean()
    fileAppendLine( G_CRONLONGTIME_FILE, ":: END Clean Sessions IP: " + dateGetNow() )
}

//Link same images files

func cronImagesLink() {
    fileAppendLine( G_CRONLONGTIME_FILE, "::Clean Duply Images IP: " + dateGetNow() )
    fileAppendLine( G_CRONLONGTIME_FILE, "" )
    datafiles := make(map[string]string)
    files := getFiles( G_IMAGES_FOLDER, "" )
    for _, file := range files {
        thash := fileHash(file)
        if _, ok := datafiles[thash]; ok {
            if datafiles[thash] != file {
                os.Remove( file )
                os.Link( datafiles[thash], file )
            }
        } else {
            datafiles[thash] = file
        }
    }
    fileAppendLine( G_CRONLONGTIME_FILE, ":: END Duply Images IP: " + dateGetNow() )
}

//Clean not ident files

func cronCleanMediaNotIdent() {
    fileAppendLine( G_CRONLONGTIME_FILE, "::Clean Not Ident Files: " + dateGetNow() )
    fileAppendLine( G_CRONLONGTIME_FILE, "" )
    media := sqlite_getMediaIdentBad( 1000 )
    tcomp := time.Now().Add(time.Hour * time.Duration((24 * G_DOWN_SAFEDAYS * -1)))
    for _, m := range media {
        if fileExist(m["file"]) && fileModifTime(m["file"]).Before(tcomp) {
            sqlite_media_delete(m["idmedia"])
            fileRemove(m["file"])
            fileAppendLine( G_CRONLONGTIME_FILE, "-- Delete: " + filepath.Base(m["file"]) )
        }
    }
    fileAppendLine( G_CRONLONGTIME_FILE, ":: END Clean Not Ident Files: " + dateGetNow() )
}

//Clean MediaInfo Duplys

func cronCleanMediaInfoDuplys() {
    fileAppendLine( G_CRONLONGTIME_FILE, "::Clean MediaInfo Duplys: " + dateGetNow() )
    fileAppendLine( G_CRONLONGTIME_FILE, "" )
    mediainfo := sqlite_getMediaInfoDuplys()
    for _, mi := range mediainfo {
        fileAppendLine( G_CRONLONGTIME_FILE, "-- FINDED: " + mi["title"] + "/" + mi["year"] + "/" + mi["season"] + "/" + mi["episode"] + "/" )
        mediainfo2 := sqlite_getMediaInfoExist( mi["title"], mi["year"], mi["season"], mi["episode"] )
        for _, mi2 := range mediainfo2 {
            if mi2["idmediainfo"] != mi["idmediainfo"] && mi2["title"] == mi["title"] && mi2["year"] == mi["year"] && mi2["season"] == mi["season"] && mi2["episode"] == mi["episode"] {
                fileAppendLine( G_CRONLONGTIME_FILE, "-- SET: " + mi2["idmediainfo"] + " => " + mi["idmediainfo"] )
                sqlite_mediainfo_delete(mi2["idmediainfo"])
                sqlite_media_change_mediainfo(mi2["idmediainfo"], mi["idmediainfo"])
            }
        }
    }
    fileAppendLine( G_CRONLONGTIME_FILE, ":: END Clean MediaInfo Duplys: " + dateGetNow() )
}

//Clean Dirs low size on downloads folder

func cronCleanDirsLowSize() {
    fileAppendLine( G_CRONLONGTIME_FILE, "::Clean Low Size Dirs: " + dateGetNow() )
    fileAppendLine( G_CRONLONGTIME_FILE, "" )
    
    folders := getFolders(G_DOWNLOADS_FOLDER)
    for _, path := range folders {
        if fileExist( path ) && checkIsDir( path ) && sliceInString( path, G_DOWNLOADS_FOLDER_EXC ) == false && G_DOWN_DIRMINSIZE > 0 && dirSizeMB(path) <= float64(G_DOWN_DIRMINSIZE) {
            fileAppendLine( G_CRONLONGTIME_FILE, "-- Delete Folder: " + path )
            delTree(path)
        }
    }
    
    fileAppendLine( G_CRONLONGTIME_FILE, ":: END Clean Low Size Dirs: " + dateGetNow() )
}

//Complete series without images with chapter with image

func cronMediaInfoCompleteImgs() {
    fileAppendLine( G_CRONLONGTIME_FILE, "::Complete Series Images: " + dateGetNow() )
    fileAppendLine( G_CRONLONGTIME_FILE, "" )
    mediainfo := sqlite_getMediaInfoSeriesAll()
    for pos, mi := range mediainfo {
        fileposter := pathJoin(G_IMAGES_FOLDER, mi["idmediainfo"] + ".poster")
        if fileExist(fileposter) {
            for pos2, mi2 := range mediainfo {
                fileposter2 := pathJoin(G_IMAGES_FOLDER, mi2["idmediainfo"] + ".poster")
                if pos2 > pos && mi2["title"] == mi["title"] && mi2["year"] == mi["year"] && fileExist(fileposter2) == false {
                    copyImgsMediaInfo(mi["idmediainfo"], mi2["idmediainfo"])
                    fileAppendLine( G_CRONLONGTIME_FILE, "++ Link Images: " + mi["title"] + " (" + mi["year"] + ") : " + mi["idmediainfo"] + " => " + mi2["idmediainfo"] )
                }
            }
        }
    }
    fileAppendLine( G_CRONLONGTIME_FILE, ":: END Complete Series Images: " + dateGetNow() )
}

//Clean media files duplicated idmediainfo by quality

func cronCleanMediaDuplys() {
    fileAppendLine( G_CRONLONGTIME_FILE, "::Clean Media Files Duplys: " + dateGetNow() )
    fileAppendLine( G_CRONLONGTIME_FILE, "" )
    media := sqlite_getMediaDuplys()
    for _, mi := range media {
        fileAppendLine( G_CRONLONGTIME_FILE, "-- FINDED: " + mi["idmedia"] + "/" + mi["idmediainfo"] )
        media2 := sqlite_getMediaWhithMediaInfoID( mi["idmediainfo"] )
        if len(media2) > 1 {
            idbestfile := ""
            bestfile := ""
            for _, mi2 := range media2 {
                if bestfile == "" {
                    //first
                    bestfile = mi2["file"]
                    idbestfile = mi2["idmedia"]
                } else {
                    fileAppendLine( G_CRONLONGTIME_FILE, "-- COMPARE: " + mi2["file"] + " => " + bestfile )
                    wfile := fscrap_getWorstVideo( bestfile, mi2["file"] )
                    fileAppendLine( G_CRONLONGTIME_FILE, "-- WORSTFILE: " + wfile )
                    if wfile == bestfile {
                        tcomp := time.Now().Add(time.Hour * time.Duration((24 * G_DOWN_SAFEDAYS * -1)))
                        if fileExist(bestfile) && fileModifTime(bestfile).Before(tcomp) {
                            sqlite_media_delete( idbestfile )
                            fileRemove( bestfile )
                        }
                        bestfile = mi2["file"]
                        idbestfile = mi2["idmedia"]
                    } else if wfile == mi2["file"] {
                        tcomp := time.Now().Add(time.Hour * time.Duration((24 * G_DOWN_SAFEDAYS * -1)))
                        if fileExist(mi2["file"]) && fileModifTime(mi2["file"]).Before(tcomp) {
                            sqlite_media_delete( mi2["idmedia"] )
                            fileRemove( mi2["file"] )
                        }
                    }
                }
            }
            fileAppendLine( G_CRONLONGTIME_FILE, "-- BEST FILE: " + idbestfile + " => " + bestfile )
        }
    }
    fileAppendLine( G_CRONLONGTIME_FILE, ":: END Clean Media Files Duplys: " + dateGetNow() )
}

//Cron Uncompress Files to file_ folder

func cronUncompressFiles() {
    //G_DOWNLOADS_FOLDER
    fileAppendLine( G_CRONLONGTIME_FILE, "::Uncompress New Files: " + dateGetNow() )
    nowpath := pathAbs(G_DOWNLOADS_FOLDER)
    fileAppendLine( G_CRONLONGTIME_FILE, "::Folder: " + nowpath )
    fileAppendLine( G_CRONLONGTIME_FILE, "" )
    err := filepath.Walk(nowpath,
        func(path string, info os.FileInfo, err error) error {
            if err != nil {
                return err
            }
            if err == nil && fileExist( path ) && checkIsFile( path ) && sliceInString( path, G_DOWNLOADS_FOLDER_EXC ) == false && checkMimeCompress( path ) && fileExist( path + "_" ) == false {
                showInfo( "CRON-LONG-UNCOMPRESS: " + path )
                if strEndWith( path, ".zip" ) && Unzip( path, path + "_" ) {
                    showInfo( "CRON-LONG-UNCOMPRESSED-ZIP: " + path )
                } else if strEndWith( path, ".rar" ) && Unrar( path, path + "_" ) {
                    showInfo( "CRON-LONG-UNCOMPRESSED-RAR: " + path )
                } else if Un7z( path, path + "_" ) {
                    showInfo( "CRON-LONG-UNCOMPRESSED-7Z: " + path )
                } else {
                    showInfo( "CRON-LONG-UNCOMPRESS-ERROR: " + path )
                }
            }
            return nil
    })
    fileAppendLine( G_CRONLONGTIME_FILE, ":: END Uncompress New Files: " + dateGetNow() )
    if err != nil {
        showInfoError(err)
    }
}

//Cron clean disk if G_DOWNLOADS_FOLDER low size < G_DOWN_LOWSPACE with G_DOWN_CLEANMODE

func cronCleanDisk() {
    fileAppendLine( G_CRONLONGTIME_FILE, "::Clean LowDiskSpace: " + dateGetNow() )
    nowpath := pathAbs(G_DOWNLOADS_FOLDER)
    fileAppendLine( G_CRONLONGTIME_FILE, "::Folder: " + nowpath )
    fileAppendLine( G_CRONLONGTIME_FILE, "" )
    disksize := getFreeSpace( nowpath )
    fileAppendLine( G_CRONLONGTIME_FILE, ":: FREE DISK: " + sizeToHuman(disksize) )
    if disksize < int64( G_DOWN_LOWSPACE * 1024 * 1024 * 1024 ) {
        fileAppendLine( G_CRONLONGTIME_FILE, ":: LOWDISKSPACE: " + sizeToHuman(disksize) + " Mode: " + G_DOWN_CLEANMODE )
        switch G_DOWN_CLEANMODE {
            case "releaseold":
                cronCleanDiskReleaseOld()
            case "bigsize":
                cronCleanDiskBigSize()
            case "old":
                fallthrough
            default:
                //old
                cronCleanDiskOld()
        }
    }
    fileAppendLine( G_CRONLONGTIME_FILE, ":: END Clean LowDiskSpace: " + dateGetNow() )
}

func cronCleanDiskOld() {
    media := sqlite_getMediaAllRev()
    quantity := 100
    nowpath := pathAbs(G_DOWNLOADS_FOLDER)
    var keys []int
    for k := range media {
        keys = append(keys, k)
    }
    sort.Ints(keys)
    for _, k := range keys {
        m := media[k]
        if fileExist(m["file"]) {
            fsize := fileSize(m["file"])
            fileAppendLine( G_CRONLONGTIME_FILE, ":: Remove: " + m["idmedia"] + " => " + sizeToHuman(fsize) + " => " + filepath.Base(m["file"]) )
            fileRemove(m["file"])
            sqlite_media_delete(m["idmedia"])
        } else {
            //fileAppendLine( G_CRONLONGTIME_FILE, ":: NOT EXIST: " + m["idmedia"] )
            //autoremove
            sqlite_media_delete(m["idmedia"])
        }
        disksize := getFreeSpace( nowpath )
        quantity--
        if quantity <= 0 || disksize > int64( G_DOWN_LOWSPACE * 1024 * 1024 * 1024 ) {
            break
        }
    }
}

func cronCleanDiskReleaseOld() {
    media := sqlite_getMediaMediaInfoYear()
    quantity := 100
    nowpath := pathAbs(G_DOWNLOADS_FOLDER)
    var keys []int
    for k := range media {
        keys = append(keys, k)
    }
    sort.Ints(keys)
    for _, k := range keys {
        m := media[k]
        if fileExist(m["file"]) {
            fsize := fileSize(m["file"])
            fileAppendLine( G_CRONLONGTIME_FILE, ":: Remove: " + m["idmedia"] + " => " + sizeToHuman(fsize) + " => " + filepath.Base(m["file"]) )
            fileRemove(m["file"])
            sqlite_media_delete(m["idmedia"])
        } else {
            //fileAppendLine( G_CRONLONGTIME_FILE, ":: NOT EXIST: " + m["idmedia"] )
            //autoremove
            sqlite_media_delete(m["idmedia"])
        }
        disksize := getFreeSpace( nowpath )
        quantity--
        if quantity <= 0 || disksize > int64( G_DOWN_LOWSPACE * 1024 * 1024 * 1024 ) {
            break
        }
    }
}

func cronCleanDiskBigSize() {
    media := sqlite_getMediaAllRev()
    quantity := 1000000
    maxsize := G_DOWN_SIZEPRIO * 2
    nowpath := pathAbs(G_DOWNLOADS_FOLDER)
    var keys []int
    for k := range media {
        keys = append(keys, k)
    }
    sort.Ints(keys)
    for nowsize := maxsize; nowsize > 1024; nowsize = nowsize - 128 {
        for _, k := range keys {
            m := media[k]
            //fileAppendLine( G_CRONLONGTIME_FILE, ":: Check Size: " + sizeToHuman( int64( nowsize * 1024 * 1024 ) ) )
            if fileExist(m["file"]) {
                fsize := fileSize(m["file"])
                if fsize > int64( nowsize * 1024 * 1024 ) {
                    fileAppendLine( G_CRONLONGTIME_FILE, ":: Remove: " + m["idmedia"] + " => " + sizeToHuman(fsize) + " => " + filepath.Base(m["file"]) + " => " + sizeToHuman( int64( nowsize * 1024 * 1024 ) ) )
                    fileRemove(m["file"])
                    sqlite_media_delete(m["idmedia"])
                }
            } else {
                //fileAppendLine( G_CRONLONGTIME_FILE, ":: NOT EXIST: " + m["idmedia"] )
                //autoremove
                sqlite_media_delete(m["idmedia"])
            }
            disksize := getFreeSpace( nowpath )
            quantity--
            if quantity <= 0 || disksize > int64( G_DOWN_LOWSPACE * 1024 * 1024 * 1024 ) {
                nowsize = 0
                break
            }
        }
    }
}

//CRON LIVETV CLEAN AND UPDATE

func cronLiveTV() {
    last := sqlite_getLogsLiveTV()
    if len(last) == 0 {
        sqlite_log_insert( "livetv-check-all", "cron", "", "", "cron", "127.0.0.1" )
        fileAppendLine( G_CRONLONGTIME_FILE, "::Clean LiveTV: " + dateGetNow() )
        fileAppendLine( G_CRONLONGTIME_FILE, "" )
        
        medialive := sqlite_getMediaLiveAll()
        total := len(medialive)
        valid := 0
        nerror := 0
        for _, ml := range medialive {
            //action := ""
            codec := ffprobeVideoCodec( ml["url"] )
            if codec != "NOCODEC" {
                valid++
            } else {
                nerror++
                sqlite_medialive_delete( ml["idmedialive"] )
            }
        }
        fileAppendLine( G_CRONLONGTIME_FILE, "Result (valid/error/total): " + intToStr( valid ) + "/" + intToStr( nerror ) + "/" + intToStr( total ) )
        fileAppendLine( G_CRONLONGTIME_FILE, "::END Clean LiveTV: " + dateGetNow() )
        
        fileAppendLine( G_CRONLONGTIME_FILE, "::Add LiveTVURLS: " + dateGetNow() )
        medialive = sqlite_getMediaLiveURLAll()
        total = len(medialive)
        valid = 0
        nerror = 0
        etadded := 0
        etexist := 0
        etnerror := 0
        ettotal := 0
        for _, ml := range medialive {
            data := urlGet( ml["url"] )
            if len(data) > 100 {
                valid++
                //Extract data
                eadded, eexist, eerrors, etotal := liveTvDataAdd( data )
                etadded += eadded
                etexist += eexist
                etnerror += eerrors
                ettotal += etotal
                //TODO on total == 0 delete url
            } else {
                nerror++
                sqlite_medialiveurl_delete( ml["idmedialive"] )
            }
        }
        fileAppendLine( G_CRONLONGTIME_FILE, "URLS Result (valid/error/total): " + intToStr( valid ) + "/" + intToStr( nerror ) + "/" + intToStr( total ) )
        fileAppendLine( G_CRONLONGTIME_FILE, "MediaURLs Added (added/exist/errors/total): " + intToStr(etadded) + "/" + intToStr(etexist) + "/" + intToStr(etnerror) + "/" + intToStr(ettotal) )
        fileAppendLine( G_CRONLONGTIME_FILE, "::END Add LiveTVURLS: " + dateGetNow() )
    }
}