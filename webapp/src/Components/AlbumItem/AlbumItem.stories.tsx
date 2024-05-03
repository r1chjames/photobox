import {AlbumItem} from './AlbumItem';


// const photoResp = {"id":"p1","name":"photo 1","filesystemPath":"/tmp/photo1.jpg","albumId":"a1","tags":"","metadata":{"a":"","b":""},"createdAt":""}
// const countResp = {"photoCount": 1};


export default {
    component: AlbumItem,
    title: 'AlbumItem',
    tags: [''],
};

export const Default = {
    args: {
        task: {
            id: '1',
            title: 'Test Task',
            state: 'TASK_INBOX',
        },
    },
};