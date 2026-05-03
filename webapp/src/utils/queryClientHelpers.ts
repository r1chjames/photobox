import { QueryClient } from '@tanstack/react-query';
import { Photo } from '../Models/Photo';

export function optimisticallyUpdatePhoto(
  queryClient: QueryClient,
  photoId: string,
  partial: Partial<Photo>
) {
  queryClient.setQueriesData(
    { queryKey: ['albumPhotos'] },
    (oldData: any) => {
      if (!oldData) return oldData;
      return {
        ...oldData,
        pages: oldData.pages.map((page: any) => ({
          ...page,
          data: page.data.map((photo: any) =>
            photo.id === photoId ? { ...photo, ...partial } : photo
          ),
        })),
      };
    }
  );
}
