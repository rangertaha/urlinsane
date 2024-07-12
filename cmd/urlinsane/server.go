// // Copyright (C) 2024  Rangertaha <rangertaha@gmail.com>
// //
// // This program is free software: you can redistribute it and/or modify
// // it under the terms of the GNU General Public License as published by
// // the Free Software Foundation, either version 3 of the License, or
// // (at your option) any later version.
// //
// // This program is distributed in the hope that it will be useful,
// // but WITHOUT ANY WARRANTY; without even the implied warranty of
// // MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// // GNU General Public License for more details.
// //
// // You should have received a copy of the GNU General Public License
// // along with this program.  If not, see <http://www.gnu.org/licenses/>.
package cmd

// import (
// 	"github.com/spf13/cobra"

// 	"github.com/rangertaha/urlinsane/server"
// )

// // serverCmd represents the server command
// var serverCmd = &cobra.Command{
// 	Use:   "server",
// 	Short: "Start a websocket server to use this tool programmatically",
// 	Long:  `This command starts up a REST API server to use this tool programmatically.`,
// 	Run: func(cmd *cobra.Command, args []string) {
// 		host, _ := cmd.Flags().GetString("host")
// 		port, _ := cmd.Flags().GetString("port")
// 		concurrency, _ := cmd.Flags().GetInt("concurrency")
// 		server.NewWebSocketServer(host, port, concurrency)
// 	},
// }

// func init() {
// 	rootCmd.AddCommand(serverCmd)
// 	serverCmd.Flags().StringP("host", "a", "127.0.0.1", "IP address for API server")
// 	serverCmd.Flags().StringP("port", "p", "8080", "Port to use")
// 	serverCmd.Flags().IntP("concurrency", "c", 50, "Number of concurrent workers")
// }
