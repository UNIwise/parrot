import { Button, Typography } from "@mui/joy";
import { FC } from "react";

type ProjectTableRowProps = {
  projectId: number;
  projectName: string;
  createdAt: string;
  numberOfVersions: number;
}

export const ProjectTableRow: FC<ProjectTableRowProps> = ({
  projectId,
  projectName,
  numberOfVersions,
  createdAt,
}) => {
  const formatIsoDateToLocaleString = (isoDate: string) => {
    return new Date(isoDate).toLocaleString();
  };

  const createdAtDate = formatIsoDateToLocaleString(createdAt);

  return (
    <tr onClick={() => window.location.href = `/projects/${projectId}/versions`}>
      <td style={{ paddingLeft: "1.5rem", fontSize: "0.9rem", }}>
        <Typography level="body-xs" fontSize={'0.9rem'} >{projectName}</Typography>
      </td>

      <td style={{ paddingLeft: "0.5rem", }}>
        <Typography level="body-xs" fontSize={'0.9rem'} fontWeight={400}>{numberOfVersions}</Typography>
      </td>

      <td style={{ paddingLeft: "0.5rem" }}>
        <Typography level="body-xs" fontSize={'0.9rem'} fontWeight={400}>{createdAtDate}</Typography>
      </td>

      <td
        style={{
          textAlign: "end",
          padding: "0.5rem 1rem",
          verticalAlign: "center",
        }}
      >
        <Button sx={{ backgroundColor: '#0078ff' }} href="">All versions</Button>
      </td>
    </tr >

  );
};
