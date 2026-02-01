# GoMFT

A forensic tool for parsing and correlating NTFS artifacts. It reads the Master File Table (`$MFT`) and Update Sequence Number Journal (`$UsnJrnl:$J`) to reconstruct file system activity.

## Capabilities

- **MFT Parsing:** Iterates raw MFT records to build a map of File Reference Numbers (FRN) to full file paths.
    
- **Journal Streaming:** Reads the `$J` log in 64KB chunks to process historical file events.
    
- **Correlation:** Links journal events to MFT paths and resolves file rename operations (Old Name -> New Name).
    
- **Visualization:** Displays a live ASCII bar chart of event density during processing.
    
- **Output:** Writes processed events to a JSONL file.
    

## Build

You need Go installed. Run the following command from the project root:



```
go build -o gomft.exe .\cmd\main.go
```

## Usage

This tool requires extracted copies of the `$MFT` and `$UsnJrnl:$J` files. You cannot run it directly against locked system files on a live C: drive.

**Command:**



```
.\gomft.exe -mft <path_to_mft> -journal <path_to_journal> -out <output_file.jsonl>
```

**Example:**



```
.\gomft.exe -mft mft.bin -journal journal.bin -out results.jsonl
```

**Flags:**

- `-mft`: Path to the extracted `$MFT` file.
    
- `-journal`: Path to the extracted `$UsnJrnl:$J` file.
    
- `-out`: Destination path for the JSONL logs (default: `output.jsonl`).
