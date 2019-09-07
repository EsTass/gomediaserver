package main

import (
    "strings"
    "net/url"
	"os/exec"
    //"fmt"
)

//CHECK MAGNET

func p2p_checkMagnet( s string ) bool {
    result := false
    s = strings.Trim( s, " " )
    re := `(magnet:[^"'\s]*)`
    f := regExpGetDataFirst( s, re )
    if len(f) > 0 {
        result = true
    }
    
    return result
}

func p2p_getMagnetTitle( s string ) string {
    result := ""
    s = strings.Trim( s, " " )
    re := `&dn=(.*?)&`
    result = regExpGetDataFirst( s, re )
    result, _ = url.QueryUnescape(result)
    result = fscrap_cleanTitle(result)
    
    return result
}

//CHECK ELINK

func p2p_checkElink( s string ) bool {
    result := false
    s = strings.Trim( s, " " )
    re := `(ed2k\:\/\/\|file\|.{1,250}\|[0-9]{8,12}\|[0-9A-F]{32}\|\/)`
    f := regExpGetDataFirst( s, re )
    if len(f) > 0 {
        result = true
    }
    
    return result
}

func p2p_getElinkTitle( s string ) string {
    result := ""
    s = strings.Trim( s, " " )
    re := `file\|(.*?)\|\d+\|`
    result = regExpGetDataFirst( s, re )
    result = fscrap_cleanTitle(result)
    return result
}

//CHECK TORRENT FILE

func p2p_checkTorrent( file string ) bool {
    result := false
    mime := fileMime( file )
    mimes := []string { "application/x-bittorrent", "application/force-download", "application/torrent", "torrent" }
    if sliceInString( mime, mimes ) {
        result = true
    }
    
    return result
}

func p2p_getLinkTitle( s string ) string {
    result := ""
    domain := ""
    u, err := url.Parse(s)
    if err == nil {
        parts := strings.Split(u.Hostname(), ".")
        domain = parts[len(parts)-2] + "." + parts[len(parts)-1]
    } else {
        domain = "nodomain.com_"
    }
    result = domain + getRandomString( 6 )
    
    return result
}

//JDOWNLOADER

func p2p_checkJDownloader( s string ) bool {
    result := false
    s = strings.Trim( s, " " )
    if sliceInString( s, G_P2P_JDDOMAINS ) {
        result = true
    }
    
    return result
}

//EXTRACT DATA LIST

//type, title, link
func p2pExtractData(data string) map[int]map[string]string {
    result := make(map[int]map[string]string)
    datalines := strings.Split( data, "\n" )
    showInfo( ":: P2P-EXTRACT-DATALINES: " + intToStr( len(datalines) ) )
    for _, line := range datalines {
         if p2p_checkElink(line) {
            link := line
            title := p2p_getElinkTitle(link)
            typelink := "ELINK"
            l := map[string]string { "Type" : typelink, "Title" : title, "Link" : link }
            result[len(result)+1] = l
             cmdElink(line)
        } else if p2p_checkMagnet(line) {
            link := line
            title := p2p_getMagnetTitle(link)
            typelink := "MAGNET"
            l := map[string]string { "Type" : typelink, "Title" : title, "Link" : link }
            result[len(result)+1] = l
             cmdMagnet(line)
        } else if p2p_checkJDownloader(line) {
            link := line
            title := p2p_getLinkTitle(link)
            typelink := "JDOWNLOADER"
            l := map[string]string { "Type" : typelink, "Title" : title, "Link" : link }
            result[len(result)+1] = l
            cmdJDownloader(line)
        } else if isValidUrl(line) {
            link := line
            title := p2p_getLinkTitle(link)
            typelink := "URL"
            l := map[string]string { "Type" : typelink, "Title" : title, "Link" : link }
            result[len(result)+1] = l
             //generic?
            //cmdJDownloader(line)
        }
    }
    showInfo( ":: P2P-EXTRACT-DATALINES-RESULT: " + intToStr(len(result)) )
    
    return result
}

//ACTIOS FOR LINKS TYPES

func cmdElink( link string ) {
    if G_P2P_ELINKSCMD != "" && fileExist(G_P2P_ELINKSCMD) {
        fileAppendLine(G_P2P_ELINKSCMD, link)
    } else if G_P2P_ELINKSCMD != "" {
        cmd := strings.Replace( G_P2P_ELINKSCMD, `%ELINK%`, link, -1 )
        cmd = strings.Replace( cmd, "  ", " ", -1 )
        args := strings.Split(cmd, " ")
        if len(args) > 1 {
            command := args[0]
            i := 0
            copy(args[i:], args[i+1:])
            args[len(args)-1] = ""
            args = args[:len(args)-1]
            result, _ := exec.Command( command, args... ).Output()
            showInfo( string(result) )
        }
    }
}

func cmdMagnet( link string ) {
    if G_P2P_MAGNETSCMD != "" && fileExist(G_P2P_MAGNETSCMD) {
        fileAppendLine(G_P2P_MAGNETSCMD, link)
    } else if G_P2P_MAGNETSCMD != "" {
        cmd := strings.Replace( G_P2P_MAGNETSCMD, `%MAGNET%`, link, -1 )
        cmd = strings.Replace( cmd, "  ", " ", -1 )
        args := strings.Split(cmd, " ")
        if len(args) > 1 {
            command := args[0]
            i := 0
            copy(args[i:], args[i+1:])
            args[len(args)-1] = ""
            args = args[:len(args)-1]
            result, _ := exec.Command( command, args... ).Output()
            showInfo( string(result) )
        }
    }
}

func cmdJDownloader( link string ) {
    if G_P2P_JDFOLDER != "" && fileExist(G_P2P_JDFOLDER) {
        link = strings.Replace(link, "http:", "https:", -1)
        fname := p2p_getLinkTitle(link) + ".crawljob"
        file := pathJoin(G_P2P_JDFOLDER, fname)
        fileAppendLine( file, link)
    }
}