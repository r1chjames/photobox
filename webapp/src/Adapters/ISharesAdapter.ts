export interface ShareLink {
    token: string;
    resourceType: 'photo' | 'album';
    resourceId: string;
    url: string;
    expiry?: string;
    createdAt: string;
}

export interface ISharesAdapter {
    createShare(resourceType: 'photo' | 'album', resourceId: string, expiry?: string, password?: string): Promise<ShareLink>;
    getShares(): Promise<ShareLink[]>;
    revokeShare(token: string): Promise<void>;
}
