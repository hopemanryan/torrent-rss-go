
import styles from '../styles/Home.module.css'
import {useEffect, useState} from "react";

export default function Home() {
    const [files, setFiles] = useState([])
    const fetchData = async () => {
        try {
            const res = await fetch('/api/hello')
            const {files: rawFiles } = await res.json()
            setFiles(rawFiles)
        }   catch (e) {
            console.log(e)
        }
    }

    const deleteItem = async (file) => {
        try {
            fetch("/api/items/delete", {
                body: {
                    fileName: file
                },
                method: 'post'
            })
            setFiles(files.filter(x => x !== file ))

        }catch (e) {

        }
    }
    useEffect( () => {
       fetchData()
    }, [])
  return (
    <div className={styles.container}>
      List of items
        <br/><br/>
        {
            files?.map((file, index) => (
                <div key={index} className={styles.listItem} onClick={() => deleteItem(file)}>
                    {file}
                </div>
            ))
        }

    </div>
  )
}
