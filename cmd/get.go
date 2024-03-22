package cmd

import (
	"log"
	"fmt"
	"github.com/spf13/cobra"
)

type Role struct {
	ApplicationId          string      `json:"ApplicationId"`
	Name                   string      `json:"Name"`
	Description            string      `json:"Description"`
	ConditionalExpression  interface{} `json:"ConditionalExpression"`
	DynamicAssignmentId    interface{} `json:"DynamicAssignmentId"`
	RoleFunctionalAbilities interface{} `json:"RoleFunctionalAbilities"`
	Id                     string      `json:"Id"`
	LastUpdate             string      `json:"LastUpdate"`
}

var pGroupName, pRoleName, pNimbusFolderName, pUserID string
var pMyGroup bool

// listCmd represents the list command
var GetCmd = &cobra.Command{
	Use:   "get",
	Short: "get test code",
	Long:  `this is GET code`,
        Run: func(cmd *cobra.Command, args []string) {

		
		if pRoleName != "" {
			res, err := sClient.GetRole(pRoleName, false)
			if err != nil {
				log.Fatalf("ERROR: %v", err)
			}

			var roles Role
			for _, roleMap := range res {
				role := Role{
					ApplicationId:          roleMap["ApplicationId"].(string),
					Name:                   roleMap["Name"].(string),
					Description:            roleMap["Description"].(string),
					ConditionalExpression:  roleMap["ConditionalExpression"],
					DynamicAssignmentId:    roleMap["DynamicAssignmentId"],
					RoleFunctionalAbilities: roleMap["RoleFunctionalAbilities"],
					Id:                     roleMap["Id"].(string),
					LastUpdate:             roleMap["LastUpdate"].(string),
				}
				roles = append(roles, role)
			}

			sClient.PrintRoles(roles)
		} else {
			fmt.Println("No RoleName!")	
		}



        },


}

func init() {

	RootCmd.AddCommand(GetCmd)

        GetCmd.PersistentFlags().StringVarP(&pRoleName, "role", "r", "", "Role/Studio Name. Without 'C3PO - '")
        GetCmd.PersistentFlags().StringVarP(&pGroupName, "group", "g", "", "Group Name. Should be 'C3PO - Studio Name or Role Name'")
        GetCmd.PersistentFlags().StringVarP(&pNimbusFolderName, "nimbusfolder", "n", "", "Nimbus Folder Name which is kwown as application name.")
        GetCmd.PersistentFlags().StringVarP(&pUserID, "userid", "u", "", "HUBID")
	GetCmd.PersistentFlags().BoolVarP(&pMyGroup, "mygroup", "", false, "Get all of my groups where I am an approval manager or just a member")

}

// Define custom sorting function by role ID
//func SortByRoleID(roles []Role) {
//	sort.Slice(roles, func(i, j int) bool {
//		return roles[i].Name < roles[j].Name
//	})
//}

//// Print roles neatly
//func PrintRoles(roles []Role) {
//	for _, role := range roles {
//		//fmt.Println("Role ID:", role.Id)
//		fmt.Println("Name: C3PO - " + role.Name)
//		//fmt.Println("Description:", role.Description)
//		//fmt.Println("Last Update:", role.LastUpdate)
//		// Print other fields as needed
//		//fmt.Println("----------------------------------")
//	}
//}
