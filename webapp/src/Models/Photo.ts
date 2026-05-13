export interface Photo {
  id: string;
  name: string;
  filesystemPath: string;
  sourcePath: string;
  albumId: string;
  tags: string;
  metadata: Record<string, unknown>;
  createdAt: string;
  dateTaken?: string;
  thumbnailUrl: string;
  favorite?: boolean;
  fileHash?: string;
  mediaType?: 'image' | 'video';
  duration?: number;
  width?: number;
  height?: number;
}
