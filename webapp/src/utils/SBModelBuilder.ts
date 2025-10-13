import {Photo} from "../Models/Photo";
import {Album} from "../Models/Album";
// import {faker} from "@faker-js/faker";

export class SBModelBuilder {

    private photos: Photo[] = [];
    private albums: Album[] = [];

    public newEmptyAlbum(): this {
        const albumId = getRandomInt(50);
        this.albums.push(newAlbum(`${albumId}`));
        return this;
    }

    public newAlbumWithPhotos(photoCount: number): this {
        const albumId = getRandomInt(50);
        this.albums.push(newAlbum(`${albumId}`));
        this.newPhotoCollection(photoCount, `${albumId}`);
        return this;
    }

    public newPhotoCollection(photoCount: number, albumId: string): this {
        for (let p = 1; p < photoCount + 1; p++) {
            this.photos.push(newPhoto(`${p}`, `${albumId}`));
        }
        return this;
    }

    public getPhotos() {
        return this.photos;
    }

    public getAlbums() {
        return this.albums;
    }
}

const getRandomInt = (max: number)=> Math.floor(Math.random() * max);

const getRandomPhotoImage = ()=> {
    const images: string[] = [
        "files/p/E~TS1180x0~articles/3925134721/0266554465.jpeg",
        "files/p/TS1800x1200~sample_galleries/4369435632/9827023079.jpg",
        "files/p/TS1800x1200~sample_galleries/4369435632/9931531280.jpg",
        "files/p/TS1800x1200~sample_galleries/4369435632/4352589276.jpg",
        "files/p/TS1800x1200~sample_galleries/2020427155/7573953117.jpg",
        "files/p/TS1800x1200~sample_galleries/2020427155/6463140671.jpg"
    ];
    return images[getRandomInt(images.length - 1)];
}

export const newPhoto = (photoId: string, albumId: string) => {
    const randomPhotoImage = getRandomPhotoImage();
    const metadata: Record<string, any>[] = [
        {"exif": "[{\"DateTime\": \"2024-01-01T10:00.000\"}, {\"ApertureValue\": \"101/32\"}]"},
        {"camera": "A7Cii"}
    ];
    return new Photo(`p${photoId}`, `Photo ${photoId}`, `/tmp/photo${photoId}.jpg`, randomPhotoImage, randomPhotoImage, albumId, "", metadata, "2024-01-01T10:00.000", "https://4.img-dpreview.com");
}

export const newAlbum = (id: string)=> {
    return new Album(`${id}`, `Album ${id}`, `Album ${id}`, "", "");
}

