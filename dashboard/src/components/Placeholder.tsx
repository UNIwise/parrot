import { Box, Skeleton } from "@mui/joy"
import { FC } from "react"

export const Placeholder: FC = () => {
  return (
    <Box sx={{ display: 'flex', flexDirection: 'column', justifyContent: 'space-around' }}>
      <Skeleton variant="rectangular" width={1750} height={70} sx={{ mt: 10 }} />
      <Skeleton sx={{ borderRadius: '50%', width: 150, height: 150, mt: 10, ml: 100 }} />
      <Skeleton variant="rectangular" width={1750} height={70} sx={{ mt: 70 }} />
    </Box>
  )
}
