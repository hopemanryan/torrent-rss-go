// Next.js API route support: https://nextjs.org/docs/api-routes/introduction
import fs from 'fs'
const { exec } = require("child_process");



export default async function handler(req, res) {
  const files = await fs.promises.readdir("/data/tvshows").catch(e => [])
    const filteredFiles = (files || []).filter(file => !file.startsWith('.'))
    return res.status(200).json({ files: filteredFiles })

}