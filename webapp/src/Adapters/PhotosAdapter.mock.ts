// @ts-nocheck
import {fn, Mock} from '@storybook/test';
import * as actual from './PhotosAdapter';

export * from './PhotosAdapter';

export const getAllPhotosInfo: Mock = fn(actual.PhotosAdapter).mockName('getAllPhotosInfo');
export const getPhotoInfoById: Mock = fn(actual.PhotosAdapter).mockName('getPhotoInfoById');
export const getPhotosInfoInAlbum: Mock = fn(actual.PhotosAdapter).mockName('getPhotosInfoInAlbum');
export const getPhotoCountInAlbum: Mock = fn(actual.PhotosAdapter).mockName('getPhotoCountInAlbum');
export const getPhotoImage: Mock = fn(actual.PhotosAdapter).mockName('getPhotoImage');