// @ts-nocheck
import {fn, Mock} from '@storybook/test';
import * as actualPhoto from './PhotosAdapter';
import * as actualAlbum from './AlbumsAdapter';

export * from './AlbumsAdapter';

export const getAllAlbumsInfo: Mock = fn(actualAlbum.PhotosAdapter).mockName('getAllAlbumsInfo');
export const getCountOfPhotosInAlbum: Mock = fn(actualPhoto.PhotosAdapter).mockName('getCountOfPhotosInAlbum');
export const getAlbumInfoById: Mock = fn(actualPhoto.PhotosAdapter).mockName('getAlbumInfoById');