import os
def read_lines_between(lines, start_pattern, end_pattern):
    """Reads lines between two patterns in a file.

    Args:
        filename: The name of the file to read.
        start_pattern: The pattern to start reading from.
        end_pattern: The pattern to stop reading at.

    Returns:
        A list of the lines between the two patterns.
    """
    append = False
    new_lines = ''
    for line in lines:
        if line.startswith(start_pattern):
            #print("Found Start Pattern")
            append = True
            continue
        elif line.startswith(end_pattern):
            #print("Found End Pattern")
            break
        elif append == True:
            #print("Appending")
            #print(line)
            new_lines = new_lines + line
            
    return new_lines

def make_pre(input_string):
    input_string.replace('\n','</br>')
    return '<pre>{input_string}</pre>'.format(input_string=input_string)


res_mig_template = '''
## Resource

{res_description}

|  SAP Cloudfoundry Provider |Community Cloudfoundry Provider |
| -- | -- |
|  {res_sap} |{res_community} |

<br/>
'''

ds_mig_template = '''
## Datasource


{ds_description}

|  SAP Cloudfoundry Provider | Community Cloudfoundry Provider  |
| -- | -- |
| {ds_sap}|{ds_community} |  
'''

import_folder = 'import_template/'
export_folder = 'migration/'

if __name__ == '__main__':
    
    for file in os.listdir('import_template'):
         lines = ''
         if not file.endswith('.md') :
             continue
         filename = import_folder+file
         with open(filename) as f:
           lines = f.readlines()
         title = lines[0]
         # Resources
         start_pattern = '#RES.COMM'
         end_pattern = '##RES.COMM'
         res_comm = read_lines_between(lines, start_pattern, end_pattern).replace('\n','</br>')
         start_pattern = '#RES.SAP'
         end_pattern = '##RES.SAP'
         res_sap = read_lines_between(lines, start_pattern, end_pattern).replace('\n','</br>')
         start_pattern = '#RES.DESC'
         end_pattern = '##RES.DESC'
         res_desc = read_lines_between(lines, start_pattern, end_pattern)

         #Datasource
         start_pattern = '#DS.COMM'
         end_pattern = '##DS.COMM'
         ds_comm = read_lines_between(lines, start_pattern, end_pattern).replace('\n','</br>')
         start_pattern = '#DS.SAP'
         end_pattern = '##DS.SAP'
         ds_sap = read_lines_between(lines, start_pattern, end_pattern).replace('\n','</br>')
         start_pattern = '#DS.DESC'
         end_pattern = '##DS.DESC'
         ds_desc = read_lines_between(lines, start_pattern, end_pattern)
         ds_doc = ''
         res_doc = ''
         if  not (ds_comm == ''):
                ds_doc = ds_mig_template.format(ds_description=ds_desc,
                               ds_community=make_pre(ds_comm),ds_sap=make_pre(ds_sap))
         if not (res_comm == ''):
                res_doc = res_mig_template.format(res_description=res_desc,
                                res_community=make_pre(res_comm),res_sap=make_pre(res_sap))
         

         generated_doc = res_doc + ds_doc
         if not (generated_doc == ''):
             generated_doc = title + '\n' + generated_doc
             outfile = export_folder+file
             with open(outfile,'w') as o:
                o.write(generated_doc)
         

