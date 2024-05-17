// @ts-nocheck
import {fn, Mock} from '@storybook/test';
import * as actual from './PhotosAdapter';

export * from './AlbumsAdapter';

export const getAllAlbumsInfo: Mock = fn(actual.PhotosAdapter).mockName('getAllAlbumsInfo');
export const getCountOfPhotosInAlbum: Mock = fn(actual.PhotosAdapter).mockName('getCountOfPhotosInAlbum');
export const getAlbumInfoById: Mock = fn(actual.PhotosAdapter).mockName('getAlbumInfoById');