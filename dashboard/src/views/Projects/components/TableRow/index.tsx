import { Button, Typography } from "@mui/joy";
import { FC } from "react";

interface ProjectTableRowProps {
  id: number;
  name: string;
  createdAt: string;
  numberOfVersions: number;
}

export const ProjectTableRow: FC<ProjectTableRowProps> = ({
  id,
  name,
  numberOfVersions,
  createdAt,
}) => {
  const formatIsoDateToLocaleString = (isoDate: string) => {
    return new Date(isoDate).toLocaleString();
  };

  const createdAtDate = formatIsoDateToLocaleString(createdAt);

  return (
    <tr onClick={() => window.location.href = `/projects/${id}/versions`}>
      <td style={{ paddingLeft: "1.5rem" }}>
        <Typography level="body-xs">{name}</Typography>
      </td>

      <td style={{ paddingLeft: "0.5rem" }}>
        <Typography level="body-xs">{numberOfVersions}</Typography>
      </td>

      <td style={{ paddingLeft: "0.5rem" }}>
        <Typography level="body-xs">{createdAtDate}</Typography>
      </td>

      <td
        style={{
          textAlign: "end",
          padding: "0.5rem 5rem",
          verticalAlign: "center",
        }}
      >

        <Button sx={{ mb: "0.5rem", backgroundColor: '#0078ff' }} href="">See all versions</Button>
      </td>
    </tr>

  );
};
