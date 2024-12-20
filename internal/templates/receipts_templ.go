// Code generated by templ - DO NOT EDIT.

// templ: version: v0.2.778
package templates

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import templruntime "github.com/a-h/templ/runtime"

import "github.com/Light2Dark/splitpay/models"

func ReceiptTable(receipt models.Receipt) templ.Component {
	return templruntime.GeneratedTemplate(func(templ_7745c5c3_Input templruntime.GeneratedComponentInput) (templ_7745c5c3_Err error) {
		templ_7745c5c3_W, ctx := templ_7745c5c3_Input.Writer, templ_7745c5c3_Input.Context
		if templ_7745c5c3_CtxErr := ctx.Err(); templ_7745c5c3_CtxErr != nil {
			return templ_7745c5c3_CtxErr
		}
		templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templruntime.GetBuffer(templ_7745c5c3_W)
		if !templ_7745c5c3_IsBuffer {
			defer func() {
				templ_7745c5c3_BufErr := templruntime.ReleaseBuffer(templ_7745c5c3_Buffer)
				if templ_7745c5c3_Err == nil {
					templ_7745c5c3_Err = templ_7745c5c3_BufErr
				}
			}()
		}
		ctx = templ.InitializeContext(ctx)
		templ_7745c5c3_Var1 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var1 == nil {
			templ_7745c5c3_Var1 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		templ_7745c5c3_Err = templ.JSONScript("receipt", receipt).Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<script type=\"text/javascript\">\n\t\tconst receipt = JSON.parse(document.getElementById(\"receipt\").textContent)\n\n\t\tif (window.Alpine === undefined) {\n\t\t\tconsole.error(\"Alpine is not loaded\")\n\t\t\tdocument.addEventListener('alpine:init', () => {\n\t\t\t\tinitializeReceipts()\n\t\t\t})\n\t\t} else {\n\t\t\tinitializeReceipts()\n\t\t}\n\n\t\tfunction initializeReceipts() {\n\t\t\tAlpine.data('receipt', () => ({\n\t\t\t\treceipt: receipt,\n\n\t\t\t\tdeleteItem(itemID) {\n\t\t\t\t\tfor (let i = 0; i < this.receipt.Items.length; i++) {\n\t\t\t\t\t\titem = this.receipt.Items[i]\n\t\t\t\t\t\tif (item.ID === itemID) {\n\t\t\t\t\t\t\tthis.receipt.Items.splice(i, 1)\n\t\t\t\t\t\t\tthis.receipt.Subtotal = this.roundTo2DP(this.receipt.Subtotal - item.Price)\n\t\t\t\t\t\t\tthis.updateAmounts()\n\t\t\t\t\t\t\tbreak\n\t\t\t\t\t\t}\n\t\t\t\t\t}\t\n\t\t\t\t},\n\n\t\t\t\tmodifyQuantity(itemID, quantity) {\n\t\t\t\t\tfor (let i = 0; i < this.receipt.Items.length; i++) {\n\t\t\t\t\t\titem = this.receipt.Items[i]\n\t\t\t\t\t\tif (item.ID === itemID) {\n\t\t\t\t\t\t\titemPrice = item.Price / item.Quantity\n\t\t\t\t\t\t\tif (quantity < 0 && item.Quantity > 1) {\n\t\t\t\t\t\t\t\tthis.receipt.Items[i].Quantity += quantity\n\t\t\t\t\t\t\t\tthis.receipt.Items[i].Price = this.roundTo2DP(this.receipt.Items[i].Price - itemPrice)\n\t\t\t\t\t\t\t\tthis.receipt.Subtotal = this.roundTo2DP(this.receipt.Subtotal - itemPrice)\n\t\t\t\t\t\t\t\tthis.updateAmounts()\n\t\t\t\t\t\t\t} else if (quantity > 0) {\n\t\t\t\t\t\t\t\tthis.receipt.Items[i].Quantity += quantity\n\t\t\t\t\t\t\t\tthis.receipt.Items[i].Price = this.roundTo2DP(Number(this.receipt.Items[i].Price) + Number(itemPrice))\n\t\t\t\t\t\t\t\tthis.receipt.Subtotal = this.roundTo2DP(Number(this.receipt.Subtotal) + Number(itemPrice))\n\t\t\t\t\t\t\t\tthis.updateAmounts()\n\t\t\t\t\t\t\t}\n\t\t\t\t\t\t\tbreak\n\t\t\t\t\t\t}\n\t\t\t\t\t}\n\t\t\t\t},\n\n\t\t\t\tupdatePrice() {\n\t\t\t\t\tamount = 0\n\t\t\t\t\tfor (let i = 0; i < this.receipt.Items.length; i++) {\n\t\t\t\t\t\tamount += Number(this.receipt.Items[i].Price)\n\t\t\t\t\t}\n\t\t\t\t\tthis.receipt.Subtotal = this.roundTo2DP(amount)\n\t\t\t\t\tthis.updateAmounts()\n\t\t\t\t},\n\n\t\t\t\tupdateAmounts() {\n\t\t\t\t\tthis.receipt.TaxAmount = this.roundTo2DP((this.receipt.Subtotal + this.receipt.Discount) * (this.receipt.TaxPercent / 100))\n\t\t\t\t\tthis.updateTotalAmount()\n\t\t\t\t},\n\t\t\t\t\n\t\t\t\tupdateTotalAmount() {\n\t\t\t\t\tthis.receipt.TotalAmount = this.roundTo2DP(Number(this.receipt.Subtotal) + Number(this.receipt.Discount) + Number(this.receipt.TaxAmount) + Number(this.receipt.ServiceCharge))\n\t\t\t\t},\n\n\t\t\t\taddItem() {\n\t\t\t\t\tlastItemID = this.receipt.Items[this.receipt.Items.length - 1].ID\n\t\t\t\t\tthis.receipt.Items.push({\n\t\t\t\t\t\tID: lastItemID + 1,\n\t\t\t\t\t\tName: \"food\",\n\t\t\t\t\t\tQuantity: 1,\n\t\t\t\t\t\tPrice: 5\n\t\t\t\t\t})\n\t\t\t\t\tthis.updatePrice()\n\t\t\t\t},\n\n\t\t\t\thandleKeyDown(e) {\n\t\t\t\t\tif (e.key === \"Enter\") {\n\t\t\t\t\t\te.preventDefault()\n\t\t\t\t\t}\n\t\t\t\t},\n\n\t\t\t\t// TODO: is there a better way to handle this + not hardcode the redirection link\n\t\t\t\tasync saveReceipt(e) {\n\t\t\t\t\te.preventDefault()\n\t\t\t\t\t// set numeric fields to 2 decimal places\n\t\t\t\t\tthis.receipt.Subtotal = this.roundTo2DP(this.receipt.Subtotal)\n\t\t\t\t\tthis.receipt.TaxAmount = this.roundTo2DP(this.receipt.TaxAmount)\n\t\t\t\t\tthis.receipt.TotalAmount = this.roundTo2DP(this.receipt.TotalAmount)\n\t\t\t\t\tthis.receipt.ServiceCharge = this.roundTo2DP(this.receipt.ServiceCharge)\n\t\t\t\t\tfor (let i = 0; i < this.receipt.Items.length; i++) {\n\t\t\t\t\t\tthis.receipt.Items[i].Price = this.roundTo2DP(this.receipt.Items[i].Price)\n\t\t\t\t\t}\n\n\t\t\t\t\tfetch(\"/saveReceipt\", {\n\t\t\t\t\t\tmethod: \"PUT\",\n\t\t\t\t\t\theaders: {\n\t\t\t\t\t\t\t\"Content-Type\": \"application/json\"\n\t\t\t\t\t\t},\n\t\t\t\t\t\tbody: JSON.stringify(receipt)\n\t\t\t\t\t})\n\t\t\t\t\t.then(response => response.json())\n\t\t\t\t\t.then(data => {\n\t\t\t\t\t\tdocument.getElementById(\"shareLink\").textContent = window.location.href + \"viewReceipt/\" + receipt.Link\n\t\t\t\t\t\tnavigator.clipboard.writeText(document.getElementById(\"shareLink\").textContent)\n\t\t\t\t\t\talert(\"Share Link copied to clipboard!\")\n\t\t\t\t\t})\n\t\t\t\t},\n\n\t\t\t\troundTo2DP(number) {\n\t\t\t\t\treturn Number((Math.round(number * 100) / 100).toFixed(2))\n\t\t\t\t}\n\t\t\t}))\n\t\t}\n\t</script><form @submit=\"saveReceipt\" x-data=\"receipt\" @keydown=\"handleKeyDown\" class=\"w-5/6 mt-14 mb-10\"><table><thead><tr><th>Description</th><th>Qty</th><th>Price</th></tr></thead><template x-for=\"item in receipt.Items\"><tr><td><button type=\"button\" @click=\"deleteItem(item.ID)\"><i class=\"fa fa-minus-circle\"></i></button> &nbsp;&nbsp; <input type=\"text\" :name=\"`itemName-${item.ID}`\" x-model=\"item.Name\" class=\"w-5/6 rounded-sm px-0.5 text-sm\"></td><td><button type=\"button\" @click=\"modifyQuantity(item.ID, -1)\"><span class=\"text-xs\">&#9664;</span></button> <span x-text=\"item.Quantity\"></span> <button type=\"button\" @click=\"modifyQuantity(item.ID, 1)\"><span class=\"text-xs\">&#9658;</span></button></td><td><input :name=\"`itemPrice-${item.ID}`\" type=\"number\" step=\"0.01\" x-model.number=\"item.Price\" @input=\"updatePrice\" class=\"w-14 rounded-sm px-0.5 text-right\"></td></tr></template><tfoot class=\"border-t border-black font-semibold\"><tr><td colspan=\"2\">Subtotal</td><td x-text=\"receipt.Subtotal\"></td></tr><tr><td colspan=\"2\">Discount</td><td><input type=\"number\" name=\"discount\" step=\"0.01\" value=\"-0\" max=\"0\" x-model.number=\"receipt.Discount\" @input=\"updateTotalAmount\" class=\"w-14 rounded-sm px-0.5 text-right\"></td></tr><tr><td colspan=\"2\">Service Charge</td><td><input type=\"number\" name=\"service-charge\" step=\"0.01\" x-model.number=\"receipt.ServiceCharge\" @input=\"updateTotalAmount\" class=\"w-14 rounded-sm px-0.5 text-right\"></td></tr><tr><td colspan=\"2\">Tax (<span x-text=\"receipt.TaxPercent\"></span>%)</td><td x-text=\"receipt.TaxAmount\"></td></tr><tr class=\"font-extrabold border-b border-black\"><td colspan=\"2\">Total</td><td x-text=\"receipt.TotalAmount\"></td></tr></tfoot></table><div class=\"mt-10 flex flex-row gap-6 justify-center\"><button type=\"button\" @click=\"addItem\" class=\"rounded-lg border border-black p-2\">Add Item</button> <button type=\"submit\" class=\"rounded-lg border border-black p-2\">Share Link</button></div></form><p id=\"shareLink\" class=\"text-center\"></p>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		return templ_7745c5c3_Err
	})
}

//	templ Table(receipt models.Receipt) {
//		<table id="receiptTable">
//			<thead>
//				<tr>
//					<th>Description</th>
//					<th>Qty</th>
//					<th>Price</th>
//				</tr>
//			</thead>
//			for _, item := range receipt.Items {
//				<tr>
//					<td>
//						// magic
//						<button type="button" @click="deleteItem(item.ID)">
//							<i class="fa fa-minus-circle"></i>
//						</button>
//						&nbsp;&nbsp;
//						<input
//							type="text"
//							name={ fmt.Sprintf("itemName-%v", item.ID) }
//							value={ item.Name }
//							class="w-5/6 rounded-sm px-0.5 text-sm"
//						/>
//					</td>
//					<td>
//						<button type="button" @click="modifyQuantity(item.ID, -1)"><span class="text-xs">&#9664;</span></button>
//						<span>{ string(item.Quantity) }</span>
//						<button
//							type="button"
//							hx-put={ string(templ.URL(fmt.Sprintf("/modifyQuantity?itemID=%v&quantity=%d", item.ID, 1))) }
//							hx-target="#receiptTable"
//							hx-swap="outerHTML"
//						><span class="text-xs">&#9658;</span></button>
//					</td>
//					<td>
//						<input
//							name={ fmt.Sprintf("itemPrice-%v", item.ID) }
//							type="number"
//							step="0.01"
//							value={ fmt.Sprintf("%v", item.Price) }
//							class="w-14 rounded-sm px-0.5 text-right"
//						/>
//					</td>
//				</tr>
//			}
//			<tfoot class="border-t border-black font-semibold">
//				<tr>
//					<td colspan="2">Subtotal</td>
//					<td x-text="receipt.Subtotal"></td>
//				</tr>
//				<tr>
//					<td colspan="2">Discount</td>
//					<td>
//						<input
//							type="number"
//							name="discount"
//							step="0.01"
//							value="-0"
//							max="0"
//							x-model.number="receipt.Discount"
//							@input="updateTotalAmount"
//							class="w-14 rounded-sm px-0.5 text-right"
//						/>
//					</td>
//				</tr>
//				<tr>
//					<td colspan="2">Service Charge</td>
//					<td>
//						<input
//							type="number"
//							name="service-charge"
//							step="0.01"
//							x-model.number="receipt.ServiceCharge"
//							@input="updateTotalAmount"
//							class="w-14 rounded-sm px-0.5 text-right"
//						/>
//					</td>
//				</tr>
//				<tr>
//					<td colspan="2">Tax (<span x-text="receipt.TaxPercent"></span>%)</td>
//					<td x-text="receipt.TaxAmount"></td>
//				</tr>
//				<tr class="font-extrabold border-b border-black">
//					<td colspan="2">Total</td>
//					<td>{ fmt.Sprintf("%v", receipt.TotalAmount) }</td>
//				</tr>
//			</tfoot>
//		</table>
//	}
var _ = templruntime.GeneratedTemplate
