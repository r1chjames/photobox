import React, {createContext, useContext, useMemo} from 'react';
import {RestApiAdapter} from "../Adapters/RestApiAdapter";
import {PhotosAdapter} from "../Adapters/PhotosAdapter";
import {AlbumsAdapter} from "../Adapters/AlbumsAdapter";
import {SettingsAdapter} from "../Adapters/SettingsAdapter";
import {UsersAdapter} from "../Adapters/UsersAdapter";
import {SharesAdapter} from "../Adapters/SharesAdapter";
import {IPhotosAdapter} from "../Adapters/IPhotosAdapter";
import {IAlbumsAdapter} from "../Adapters/IAlbumsAdapter";
import {ISettingsAdapter} from "../Adapters/ISettingsAdapter";
import {IUsersAdapter} from "../Adapters/IUsersAdapter";
import {ISharesAdapter} from "../Adapters/ISharesAdapter";

interface AdapterContextType {
    photosAdapter: IPhotosAdapter;
    albumsAdapter: IAlbumsAdapter;
    settingsAdapter: ISettingsAdapter;
    usersAdapter: IUsersAdapter;
    sharesAdapter: ISharesAdapter;
}

const AdapterContext = createContext<AdapterContextType | null>(null);

export const AdapterProvider: React.FunctionComponent<{ baseApiUrl: string; children: React.ReactNode }> = ({baseApiUrl, children}) => {
    const adapters = useMemo(() => {
        const restApi = new RestApiAdapter(baseApiUrl);
        return {
            photosAdapter: new PhotosAdapter(restApi),
            albumsAdapter: new AlbumsAdapter(restApi),
            settingsAdapter: new SettingsAdapter(restApi),
            usersAdapter: new UsersAdapter(restApi),
            sharesAdapter: new SharesAdapter(restApi),
        };
    }, [baseApiUrl]);

    return (
        <AdapterContext.Provider value={adapters}>
            {children}
        </AdapterContext.Provider>
    );
};

export const useAdapters = (): AdapterContextType => {
    const context = useContext(AdapterContext);
    if (!context) {
        throw new Error('useAdapters must be used within an AdapterProvider');
    }
    return context;
};
