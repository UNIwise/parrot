import { Delete } from '@mui/icons-material';
import { Typography } from '@mui/joy';
import { FC } from 'react';
import { Project } from '../../../../interfaces/projects';

interface ProjectTableRowProps {
    projectInfo: Project;
}

export const ProjectTableRow: FC<ProjectTableRowProps> = ({ projectInfo }) => {
    const formatIsoDateToLocaleString = (isoDate: string) => {
        return new Date(isoDate).toLocaleString();
    };

    const createdAtDate = formatIsoDateToLocaleString(projectInfo.createdAt);

    return (
        <tr>
            <td style={{ paddingLeft: '1.5rem' }}>
                <Typography level="body-xs">{projectInfo.name}</Typography>
            </td>

            <td style={{ paddingLeft: '0.5rem' }}>
                <Typography level="body-xs">{projectInfo.numberOfVersions}</Typography>
            </td>

            <td style={{ paddingLeft: '0.5rem' }}>
                <Typography level="body-xs">{createdAtDate}</Typography>
            </td>

            <td style={{ textAlign: 'end', padding: '0.5rem 5rem', verticalAlign: 'center' }}>
                <Delete />
            </td>
        </tr>
    );
};
