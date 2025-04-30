import { Delete } from "@mui/icons-material";
import Add from "@mui/icons-material/Add";
import { FormLabel } from "@mui/joy";
import Button from "@mui/joy/Button";
import DialogContent from "@mui/joy/DialogContent";
import DialogTitle from "@mui/joy/DialogTitle";
import FormControl from "@mui/joy/FormControl";
import Input from "@mui/joy/Input";
import Modal from "@mui/joy/Modal";
import ModalDialog from "@mui/joy/ModalDialog";
import Stack from "@mui/joy/Stack";
import { ChangeEvent, FC, FormEvent, useState } from "react";
import { useDeleteVersion } from "../api/hooks/useDeleteVersion";
import { usePostVersion } from "../api/hooks/usePostVersion";

type ManageVersionModalProps = {
  projectId: number;
  versionId?: number;
  versionName?: string;
};

export const ManageVersionModal: FC<ManageVersionModalProps> = ({
  projectId,
  versionId,
  versionName,
}) => {
  const [open, setOpen] = useState(false);
  const [newProjectVersion, setNewProjectVersion] = useState("");
  const [validProjectVersion, setValidProjectVersion] = useState(false);

  const [processing, setProcessing] = useState(false);

  const { mutateAsync: deleteVersion } = useDeleteVersion(
    projectId,
    versionId!,
  );
  const { mutateAsync: postNewVersion } = usePostVersion(projectId);

  const onProjectNameChange = (event: ChangeEvent<HTMLInputElement>) => {
    const valid = /^[a-zA-Z0-9.-_]*$/.test(event.target.value);
    setValidProjectVersion(valid);

    setNewProjectVersion(event.target.value);
  };

  const onSubmit = async (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault();

    try {
      setProcessing(true);

      if (versionId) {
        await deleteVersion();
      } else {
        await postNewVersion({ name: newProjectVersion });
      }
    } finally {
      setProcessing(false);
      setOpen(false);
    }
  };

  return (
    <>
      {versionId ? (
        <Button
          variant="outlined"
          color="danger"
          sx={{ color: (t) => t.palette.danger[500] }}
          endDecorator={<Delete />}
          onClick={() => setOpen(true)}
        >
          Delete
        </Button>
      ) : (
        <Button
          color="primary"
          startDecorator={<Add />}
          onClick={() => setOpen(true)}
          sx={{
            height: "36px",
          }}
        >
          New Version
        </Button>
      )}

      <Modal open={open} onClose={() => setOpen(false)}>
        <ModalDialog sx={{ width: "400px" }}>
          <form onSubmit={onSubmit}>
            {versionId ? (
              <Stack spacing={2}>
                <DialogTitle>Delete version</DialogTitle>
                <DialogContent>
                  Are you sure you want to delete {versionName} version?
                </DialogContent>

                <Button type="submit" color="danger" loading={processing}>
                  Delete
                </Button>
              </Stack>
            ) : (
              <Stack spacing={2}>
                <DialogTitle>New version</DialogTitle>

                <DialogContent>
                  Choose a tag for the new version. Creating a new version will
                  take a couple of minutes.
                </DialogContent>

                <FormControl>
                  <FormLabel>Version tag</FormLabel>
                  <Input
                    autoFocus
                    required
                    onChange={onProjectNameChange}
                    value={newProjectVersion}
                    error={!!newProjectVersion && !validProjectVersion}
                    disabled={processing}
                  />
                </FormControl>

                <Button
                  type="submit"
                  disabled={!validProjectVersion}
                  loading={processing}
                >
                  Add version
                </Button>
              </Stack>
            )}
          </form>
        </ModalDialog>
      </Modal>
    </>
  );
};
