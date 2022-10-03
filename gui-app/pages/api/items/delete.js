import fs from "fs";

export default async function handler(req, res) {
    try {
        const {fileName} = req.body
        if (!fileName) {
            return res.status(500).json({err: "no file name"})
        }

        await fs.promises.unlink(`$/data/tvshows/${fileName}`)
        return res.status(200)

    }catch (e) {
        return res.status(500).json({err: e})
    }
}