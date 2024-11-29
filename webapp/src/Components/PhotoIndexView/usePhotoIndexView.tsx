import React, {useEffect, useState} from 'react';
import {IPhotosAdapter} from "../../Adapters/IPhotosAdapter";
import {nextPhotoInAlbum, Photo, previousPhotoInAlbum} from "../../Models/Photo";
import {PhotoItem} from "../PhotoItem/PhotoItem";

const usePhotoIndexView = (photosAdapter: IPhotosAdapter, id: string) => {

    const [photos, setPhotos] = useState<Photo[]>();
    const [photoItems, setPhotoItems] = useState<JSX.Element[]>();
    const [photosLoaded, setPhotosLoaded] = useState(false);

    const getAllPhotos = async(propsPhotosAdapter: IPhotosAdapter, albumId: string) => {
        return await propsPhotosAdapter.getPhotosInfoInAlbum(albumId, 1, 100, true);
    };

    const updatePhotoCollection = (items: Map<string, JSX.Element[]>, date: string, itemToAdd: JSX.Element) => {
        if (date == null || undefined) {
            console.log(itemToAdd)
        }
        const itemToUpdate = items.has(date) ? items.get(date) : [];
        itemToUpdate!.push(itemToAdd);
        return items.set(date, itemToUpdate!);
    };

    const loadItems = (groupKey: number, imageSources: Photo[]) => {
        let items = new Map<string, JSX.Element[]>();
        const start = 0;

        if (imageSources) {
            for (let i = 0; i < imageSources.length; i += 1) {
                const imageSource = imageSources[start + i];
                if (typeof imageSource !== 'undefined') {
                    items = updatePhotoCollection(items,
                        imageSource.getPhotoDate(),
                        (
                            <PhotoItem
                                src={imageSource.sourcePath}
                                thumbnail={imageSource.sourcePath}
                                groupKey={groupKey}
                                key={start + i}
                                source={imageSource}
                                photoSequence={1 + start + i}
                                previousPhoto={() => previousPhotoInAlbum(photos!, i)}
                                nextPhoto={() => nextPhotoInAlbum(photos!, i)}
                                lastInAlbum={i <= imageSources.length}
                            />
                        )
                    );
                }
            }
            return items;
        }
    };


    const parseItems = (items: Map<string, JSX.Element[]>) => {
        const itemsToReturn: JSX.Element[] = [];
        items.forEach((values, key) => {
            itemsToReturn.push(
                <div>
                    {key}
                    {values}
                </div>
            );
        });

        return itemsToReturn;
    };

    useEffect(() => {
        (async function retrievePhotos() {
            const retrievedPhotos = await getAllPhotos(photosAdapter, id as string);
            setPhotos(retrievedPhotos);
            const photoItems = parseItems(loadItems(0, retrievedPhotos)!);
            setPhotoItems(photoItems!);
        })().then(() => setPhotosLoaded(true));
    }, [setPhotos, setPhotoItems, id, photosAdapter]);

    return [{photoItems, photosLoaded}]
};

export default usePhotoIndexView;