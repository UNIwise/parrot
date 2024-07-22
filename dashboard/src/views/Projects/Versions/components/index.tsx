import { Delete } from "@mui/icons-material";
import { Button, Typography } from "@mui/joy";
import { FC } from "react";
import { useDeleteVersion } from "../../../../api/hooks/useDeleteVersion";
import { useGetPageParams } from "../../../../api/hooks/useGetPageParams";


interface VersionTableRowProps {
  versionId: number;
  versionName: string;
  createdAt: string;
}

export const VersionTableRow: FC<VersionTableRowProps> = ({
  versionId,
  versionName,
  createdAt,
}) => {
  const formatIsoDateToLocaleString = (isoDate: string) => {
    return new Date(isoDate).toLocaleString();
  };

  const { projectId } = useGetPageParams();
  const { mutate: deleteVersion } = useDeleteVersion(projectId, versionId);

  const createdAtDate = formatIsoDateToLocaleString(createdAt);

  return (
    <tr>
      <td style={{ paddingLeft: "1.5rem" }}>
        <Typography level="body-xs">{versionName}</Typography>
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
        <Button sx={{ mb: "0.5rem", backgroundColor: '#0078ff' }} onClick={() => deleteVersion()}>Delete <Delete /></Button>
      </td>
    </tr>
  );
};
