import {Photo} from "../Models/Photo";
import {Album} from "../Models/Album";

export const newPhoto: (id: number) => Photo = (id: number) => {
    return new Photo(`p${id}`, `Photo ${id}`, `/tmp/photo${id}.jpg`, "files/p/E~TS1180x0~articles/3925134721/0266554465.jpeg", "a1", "", {"a":"","b":""}, "", "https://4.img-dpreview.com");
}

export const newAlbum: (id: number) => Album = (id: number) => {
    return new Album(`a${id}`, `Album ${id}`, `Album ${id}`, "", "");
}

export const newAlbumWithPhotos: (id: number) => Album = (id: number) => {
    return new Album(`a${id}`, `Album ${id}`, `Album ${id}`, "", "");
}



