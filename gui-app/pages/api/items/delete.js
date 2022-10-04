import fs from "fs";

export default async function handler(req, res) {
    try {
        const {fileName} = JSON.parse(req.body)
        if (!fileName) {
            return res.status(500).json({err: "no file name"})
        }

        const path = "/data/tvshows/"+ fileName
        const pathState  = await fs.promises.stat(path)
         fs.rmSync(`/data/tvshows/${fileName}`,  { recursive: true, force: true })
        return res.status(200).json({"ok": true})

    }catch (e) {
        return res.status(500).json({err: e})
    }
}
