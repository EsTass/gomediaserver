package main


import (
	//"fmt"
    "strings"
    "encoding/json"
	"os/exec"
    //"net/url"
)

//GET URL ID (imdb,thetvdb, omdb, filmaffinity)

func wscrap_getURLID( url string ) string {
    result := ""
    
    //IMDB
    result = regExpGetDataFirst( url, `(tt\d{6,10})` )
    
    //themoviedb
    if result == "" {
        result = regExpGetDataFirst( url, `(\/movie\/[0-9]{1,10})` )
    }
    
    //thetvdb
    if result == "" {
        result = regExpGetDataFirst( url, `(id=[0-9]{1,10})` )
    }
    
    //FILMAFFINITY
    if result == "" {
        result = regExpGetDataFirst( url, `(\/film([0-9]{5,10})\.html)` )
    }
    
    return result
}

//CLEAN TITLES from websearch

func wscrap_cleanTitle( title string ) string {
    result := title
    
    result = strings.Replace( result, "- IMDb", "", -1 )
    result = strings.Replace( result, "- FilmAffinity", "", -1 )
    
    return result
}

//WEBSEARCH

func webSearch( s string, filterurl string, titleurl string ) map[string]string {
    result := make(map[string]string)
    
    d := webSearchDDG(s, filterurl, titleurl)
    if len(d) > 0 {
        mapsJoinSS(result, d)
        //result = d
    }
    if len(d) < 5 {
        d := webSearchGoogler(s, filterurl, titleurl)
        if len(d) > 0 {
            mapsJoinSS(result, d)
        }
    }
    
    return result
}

//DUCKDUCKGO

type wsDDGRTemplate []struct {
	Abstract string `json:"abstract"`
	Title    string `json:"title"`
	URL      string `json:"url"`
}

func webSearchDDG( s string, filterurl string, titleurl string ) map[string]string {
    //title => url
    result := make(map[string]string)
    //escapedQuery := url.QueryEscape(s)
    
    //over G_WS_DDG
    command := G_WS_DDG
    args := []string{
        "--json", 
        s,
    }
    
    out, _ := exec.Command( command, args... ).Output()
    //out, err := exec.Command( command, args... ).Output()
    //if err != nil {
        //log.Fatal(err)
    //}
    showInfo( "DDG-RESPONSE-SIZE: " + intToStr(len(out)) )
    //fmt.Printf("The date is %s\n", out)
    responseString := string(out)
    data := wsDDGRTemplate{}
    _ = json.Unmarshal([]byte(responseString), &data)
    for _, d := range data {
        showInfo( "DDG-RESPONSE-DATALINE: " + d.Title + " => " + d.URL )
        if ( len(filterurl) == 0 || strInStr( d.URL, filterurl ) ) && ( len(titleurl) == 0 || strInStr( d.Title, titleurl ) ) {
            result[wscrap_cleanTitle(d.Title)] = d.URL
        }
    }
    
    return result
}

//GOOGLER

func webSearchGoogler( s string, filterurl string, titleurl string ) map[string]string {
    //title => url
    result := make(map[string]string)
    //escapedQuery := url.QueryEscape(s)
    
    //over G_WS_GOOGLER
    command := G_WS_GOOGLER
    args := []string{
        "--json", 
        s,
    }
    
    out, _ := exec.Command( command, args... ).Output()
    //out, err := exec.Command( command, args... ).Output()
    //if err != nil {
        //log.Fatal(err)
    //}
    showInfo( "GOOGLER-RESPONSE-SIZE: " + intToStr(len(out)) )
    //fmt.Printf("The date is %s\n", out)
    responseString := string(out)
    data := wsDDGRTemplate{}
    _ = json.Unmarshal([]byte(responseString), &data)
    for _, d := range data {
        showInfo( "GOOGLER-RESPONSE-DATALINE: " + d.Title + " => " + d.URL )
        if ( len(filterurl) == 0 || strInStr( d.URL, filterurl ) ) && ( len(titleurl) == 0 || strInStr( d.Title, titleurl ) ) {
            result[d.Title] = d.URL
        }
    }
    
    return result
}