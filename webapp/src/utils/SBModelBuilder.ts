import {Photo} from "../Models/Photo";
import {Album} from "../Models/Album";

export class SBModelBuilder {

    private photos: Photo[] = [];
    private albums: Album[] = [];

    public newAlbumWithPhotos(albumCount: number, photoCount: number): this {
        for (let a = 1; a < albumCount + 1; a++) {
            this.albums.push(newAlbum(`${a}`));
            for (let p = 1; p < photoCount + 1; p++) {
                this.photos.push(newPhoto(`${p}`, `${a}`));
            }
        }
        return this;
    }

    public getPhotos() {
        return this.photos;
    }

    public getPhotosInAlbum(albumId: string) {
        return this.photos.filter(p => p.albumId === albumId);
    }

    public getAlbums() {
        return this.albums;
    }

}

export const newPhoto = (photoId: string, albumId: string) => {
    return new Photo(`p${photoId}`, `Photo ${photoId}`, `/tmp/photo${photoId}.jpg`, "files/p/E~TS1180x0~articles/3925134721/0266554465.jpeg", albumId, "", {"exif": {"DateTime": "2024-01-01T10:00.000"}}, "2024-01-01T10:00.000", "https://4.img-dpreview.com");
}

export const newAlbum = (id: string)=> {
    return new Album(`a${id}`, `Album ${id}`, `Album ${id}`, "", "");
}

