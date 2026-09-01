import {Photo} from "../Models/Photo";
import {Album} from "../Models/Album";

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
        "https://1.img-dpreview.com/files/p/E~TS1180x0~articles/3925134721/0266554465.jpeg",
        "https://2.img-dpreview.com/files/p/TS1800x1200~sample_galleries/4369435632/9827023079.jpg",
        "https://3.img-dpreview.com/files/p/TS1800x1200~sample_galleries/4369435632/9931531280.jpg",
        "https://4.img-dpreview.com/files/p/TS1800x1200~sample_galleries/4369435632/4352589276.jpg",
        "https://1.img-dpreview.com/files/p/TS1800x1200~sample_galleries/2020427155/7573953117.jpg",
        "https://2.img-dpreview.com/files/p/TS1800x1200~sample_galleries/2020427155/6463140671.jpg"
    ];
    return images[getRandomInt(images.length - 1)];
}

export const newPhoto = (photoId: string, albumId: string): Photo => {
    const randomPhotoImage = getRandomPhotoImage();
    const metadata: Record<string, unknown> = {
        "exif": "[{\"DateTime\": \"2024-01-01T10:00.000\"}, {\"ApertureValue\": \"101/32\"}]",
        "camera": "A7Cii"
    };
    // Varied native dimensions so the masonry grid shows mixed tile sizes:
    // ~60% landscape, ~20% portrait, ~20% square — mirroring a real library.
    const layouts: Array<[number, number]> = [
        [1920, 1280], [1280, 1920], [1920, 1080], [1080, 1080],
        [2400, 1600], [1600, 2400], [1920, 1280], [1080, 1920],
        [2048, 1536], [1200, 1200],
    ];
    const [width, height] = layouts[getRandomInt(layouts.length - 1)];
    return {
        id: `p${photoId}`,
        name: `Photo ${photoId}`,
        filesystemPath: `/tmp/photo${photoId}.jpg`,
        sourcePath: randomPhotoImage,
        albumId,
        tags: "",
        metadata,
        createdAt: "2024-01-01T10:00.000",
        thumbnailUrl: randomPhotoImage,
        width,
        height,
    };
}

export const newAlbum = (id: string): Album => {
    return { id: `${id}`, name: `Album ${id}`, description: `Album ${id}`, tags: "", metadata: "" };
}

