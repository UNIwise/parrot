import { useParams } from 'react-router-dom';

export const useGetPageParams = () => {
  const { projectId } = useParams<{
    projectId: string;
  }>();

  return {
    projectId: projectId ? parseInt(projectId) : undefined,
  };
};
