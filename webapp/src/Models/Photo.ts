export interface Photo {
  id: string;
  name: string;
  filesystemPath: string;
  sourcePath: string;
  albumId: string;
  tags: string;
  metadata: Record<string, unknown>;
  createdAt: string;
  thumbnail: string;
}
